package gogencommand

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/printer"
	"go/token"
	"io/ioutil"
	"os"
	"strings"

	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
	"golang.org/x/tools/imports"
)

type RepositoryModel struct {
	RepositoryName string
	EntityName     string
	UsecaseName    string
	PackagePath    string
}

func NewRepositoryModel() (Commander, error) {
	flag.Parse()

	// create a repository need to have 2 mandatory params and 1 optional params
	values := flag.Args()[1:]
	if len(values) < 2 {
		return nil, fmt.Errorf("repository format must include `gogen repository RepositoryName EntityName [UsecaseName]`")
	}

	usecaseName := ""
	if len(values) >= 3 {
		usecaseName = values[2]
	}

	return &RepositoryModel{
		RepositoryName: values[0],
		EntityName:     values[1],
		UsecaseName:    usecaseName,
		PackagePath:    util.GetGoMod(),
	}, nil
}

func (obj *RepositoryModel) Run() error {

	// create folder repository
	err := util.CreateFolderIfNotExist("domain/repository")
	if err != nil {
		return err
	}

	// try to create an entity
	{
		em := EntityModel{EntityName: obj.EntityName}
		err = em.Run()
		if err != nil {
			return err
		}
	}

	existingFile := fmt.Sprintf("domain/repository/repository.go")

	// create repository.go if not exist yet
	if !util.IsExist(existingFile) {
		err = util.WriteFile(templates.RepositoryFile, existingFile, obj)
		if err != nil {
			return err
		}

	} else

	// repository.go is already exist. check existing repo interface
	{

		fset := token.NewFileSet()

		pkgs, err := parser.ParseDir(fset, "domain/repository", nil, parser.ParseComments)
		if err != nil {
			if err != nil {
				return err
			}
		}

		for _, pkg := range pkgs {
			for _, file := range pkg.Files {

				for _, decl := range file.Decls {

					gen, ok := decl.(*ast.GenDecl)
					if !ok {
						continue
					}

					if gen.Tok != token.TYPE {
						continue
					}

					for _, specs := range gen.Specs {

						ts, ok := specs.(*ast.TypeSpec)
						if !ok {
							continue
						}

						if _, ok = ts.Type.(*ast.InterfaceType); !ok {
							continue
						}

						// repo already exist, abort the command
						if ts.Name.String() == fmt.Sprintf("%sRepo", obj.RepositoryName) {
							return fmt.Errorf("repo %s already exist", obj.RepositoryName)
						}
					}
				}

			}
		}
	}

	// reopen the file
	file, err := os.Open(existingFile)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	if err := file.Close(); err != nil {
		return err
	}

	// check the prefix and give specific template for it
	constTemplateCode, err := obj.prepareRepoTemplate()
	if err != nil {
		return err
	}

	// write the template in the end of file
	buffer.WriteString(constTemplateCode)
	buffer.WriteString("\n")

	// reformat and do import
	newBytes, err := imports.Process(existingFile, buffer.Bytes(), nil)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(existingFile, newBytes, 0644); err != nil {
		return err
	}

	// try to inject to outport
	obj.injectToOutport()

	return nil

}

func (obj *RepositoryModel) injectToOutport() error {

	// if the third params is not given, then nothing to do
	if obj.UsecaseName == "" {
		return nil
	}

	// read the current outport.go if it not exist we will create it
	fileReadPath := fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(obj.UsecaseName))

	fset := token.NewFileSet()

	file, err := parser.ParseFile(fset, fileReadPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	isAlreadyInjectedBefore := false

	for _, decl := range file.Decls {

		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if gen.Tok != token.TYPE {
			continue
		}

		for _, specs := range gen.Specs {

			ts, ok := specs.(*ast.TypeSpec)
			if !ok {
				continue
			}

			iFace, ok := ts.Type.(*ast.InterfaceType)
			if !ok {
				continue
			}

			// find the specific outport interface
			if ts.Name.String() != fmt.Sprintf("%sOutport", obj.UsecaseName) {
				continue
			}

			for _, meth := range iFace.Methods.List {

				selType, ok := meth.Type.(*ast.SelectorExpr)

				// if interface already injected then abort the mission
				if ok && selType.Sel.String() == fmt.Sprintf("%sRepo", obj.RepositoryName) {
					isAlreadyInjectedBefore = true
					break
				}

			}

			if !isAlreadyInjectedBefore {
				// add new repository to outport interface
				iFace.Methods.List = append(iFace.Methods.List, &ast.Field{
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "repository",
						},
						Sel: &ast.Ident{
							Name: fmt.Sprintf("%sRepo", obj.RepositoryName),
						},
					},
				})

				// after injection no need to check anymore
				break
			}

		}
	}

	if !isAlreadyInjectedBefore {

		// rewrite the outport
		f, err := os.Create(fileReadPath)
		if err := printer.Fprint(f, fset, file); err != nil {
			return err
		}
		err = f.Close()
		if err != nil {
			return err
		}

		// reformat and import
		newBytes, err := imports.Process(fileReadPath, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(fileReadPath, newBytes, 0644); err != nil {
			return err
		}

	}

	// try to inject to interactor
	obj.injectToInteractor()

	return nil

}

const injectedCodeLocation = "//!"

func (obj *RepositoryModel) injectToInteractor() error {

	existingFile := fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(obj.UsecaseName))

	// open interactor file
	file, err := os.Open(existingFile)
	if err != nil {
		return err
	}

	// check the repo name and return specific template
	constTemplateCode, err := obj.prepareInteractorTemplate()
	if err != nil {
		return err
	}

	needToInject := false

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		// check the injected code in interactor
		if strings.TrimSpace(row) == injectedCodeLocation {

			needToInject = true

			// we need to provide an error
			InitiateError()

			// inject code
			buffer.WriteString(constTemplateCode)
			buffer.WriteString("\n")

			continue
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// if no injected marker found, then abort the next step
	if !needToInject {
		return nil
	}

	if err := file.Close(); err != nil {
		return err
	}

	// rewrite the file
	if err := ioutil.WriteFile(existingFile, buffer.Bytes(), 0644); err != nil {
		return err
	}

	fset := token.NewFileSet()
	node, err := parser.ParseFile(fset, existingFile, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// reformat the code
	var newBuf bytes.Buffer
	err = format.Node(&newBuf, fset, node)
	if err != nil {
		return err
	}

	// rewrite again
	if err := ioutil.WriteFile(existingFile, newBuf.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

func (obj *RepositoryModel) prepareInteractorTemplate() (constTemplateCode string, err error) {

	if obj.hasPrefix("findone") || obj.hasPrefix("getone") {
		constTemplateCode, err = util.PrintTemplate(templates.RepositoryInjectFindOneFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	if obj.hasPrefix("find") || obj.hasPrefix("get") {
		constTemplateCode, err = util.PrintTemplate(templates.RepositoryInjectFindFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	if obj.hasPrefix("remove") || obj.hasPrefix("delete") {
		constTemplateCode, err = util.PrintTemplate(templates.RepositoryInjectRemoveFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	if obj.hasPrefix("save") || obj.hasPrefix("create") || obj.hasPrefix("add") || obj.hasPrefix("update") {
		constTemplateCode, err = util.PrintTemplate(templates.RepositoryInjectSaveFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	{
		constTemplateCode, err = util.PrintTemplate(templates.RepositoryInjectInterfaceFile, obj)
		if err != nil {
			return "", err
		}
	}

	return constTemplateCode, nil
}

func (obj *RepositoryModel) prepareRepoTemplate() (templateCode string, err error) {

	if obj.hasPrefix("findone") || obj.hasPrefix("getone") {
		templateCode, err = util.PrintTemplate(templates.RepositoryFindOneFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	if obj.hasPrefix("find") || obj.hasPrefix("get") {
		templateCode, err = util.PrintTemplate(templates.RepositoryFindFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	if obj.hasPrefix("remove") || obj.hasPrefix("delete") {
		templateCode, err = util.PrintTemplate(templates.RepositoryRemoveFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	if obj.hasPrefix("save") || obj.hasPrefix("create") || obj.hasPrefix("add") || obj.hasPrefix("update") {
		templateCode, err = util.PrintTemplate(templates.RepositorySaveFile, obj)
		if err != nil {
			return "", err
		}

	} else //

	{
		templateCode, err = util.PrintTemplate(templates.RepositoryInterfaceFile, obj)
		if err != nil {
			return "", err
		}
	}
	return templateCode, nil
}

func (obj *RepositoryModel) hasPrefix(str string) bool {
	return strings.HasPrefix(strings.ToLower(obj.RepositoryName), str)
}

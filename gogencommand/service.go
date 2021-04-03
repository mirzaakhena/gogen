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

type ServiceModel struct {
	ServiceName string
	UsecaseName string
}

func NewServiceModel() (Commander, error) {
	flag.Parse()

	values := flag.Args()[1:]
	if len(values) < 1 {
		return nil, fmt.Errorf("service format must include `gogen service ServiceName [UsecaseName]`")
	}

	usecaseName := ""
	if len(values) >= 2 {
		usecaseName = values[1]
	}

	return &ServiceModel{
		ServiceName: values[0],
		UsecaseName: usecaseName,
	}, nil
}

func (obj *ServiceModel) Run() error {

	err := util.CreateFolderIfNotExist("domain/service")
	if err != nil {
		return err
	}

	existingFile := fmt.Sprintf("domain/service/service.go")

	// create service.go if not exist yet
	if !util.IsExist(existingFile) {
		err = util.WriteFile(templates.ServiceFile, existingFile, obj)
		if err != nil {
			return err
		}

	} else

	// service.go is already exist. check existing service interface
	{
		fset := token.NewFileSet()

		pkgs, err := parser.ParseDir(fset, "domain/service", nil, parser.ParseComments)
		if err != nil {
			return err
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
						if ts.Name.String() == fmt.Sprintf("%sRepo", obj.ServiceName) {
							return fmt.Errorf("repo %s already exist", obj.ServiceName)
						}
					}
				}

			}
		}
	}

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

	constTemplateCode, err := util.PrintTemplate(templates.ServiceTemplateFile, obj)
	if err != nil {
		return err
	}

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

	err = obj.injectToOutport()
	if err != nil {
		return err
	}

	return nil

}

func (obj *ServiceModel) injectToOutport() error {

	// if the third params is not given, then nothing to do
	if obj.UsecaseName == "" {
		return nil
	}

	// read the current outport.go if it not exist we will create it
	fileReadPath := fmt.Sprintf("usecase/%s/outport.go", strings.ToLower(obj.UsecaseName))

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
			if ts.Name.String() != "Outport" {
				continue
			}

			for _, meth := range iFace.Methods.List {

				selType, ok := meth.Type.(*ast.SelectorExpr)

				// if interface already injected then abort the mission
				if ok && selType.Sel.String() == fmt.Sprintf("%sService", obj.ServiceName) {
					isAlreadyInjectedBefore = true
					break
				}

			}

			if !isAlreadyInjectedBefore {
				// add new service to outport interface
				iFace.Methods.List = append(iFace.Methods.List, &ast.Field{
					Type: &ast.SelectorExpr{
						X: &ast.Ident{
							Name: "service",
						},
						Sel: &ast.Ident{
							Name: fmt.Sprintf("%sService", obj.ServiceName),
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
		if f != nil {
			err = f.Close()
			if err != nil {
				return err
			}
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
	err = obj.injectToInteractor()
	if err != nil {
		return err
	}

	return nil

}

func (obj *ServiceModel) injectToInteractor() error {

	existingFile := fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(obj.UsecaseName))

	// open interactor file
	file, err := os.Open(existingFile)
	if err != nil {
		return err
	}

	// check the service name and return specific template
	constTemplateCode, err := util.PrintTemplate(templates.ServiceInjectFile, obj)
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
			err = InitiateError()
			if err != nil {
				return err
			}

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

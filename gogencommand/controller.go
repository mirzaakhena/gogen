package gogencommand

import (
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"strings"

	"github.com/mirzaakhena/gogen2/templates"
	"github.com/mirzaakhena/gogen2/util"
	"golang.org/x/tools/imports"
)

type ControllerModel struct {
	PackagePath      string
	UsecaseName      string
	InportMethodName string
	ControllerName   string
}

func NewControllerModel() (Commander, error) {
	flag.Parse()

	// controller need only one parameter
	values := flag.Args()[1:]
	if len(values) != 2 {
		return nil, fmt.Errorf("controller name must not empty. `gogen controller ControllerName UsecaseName`")
	}

	return &ControllerModel{
		PackagePath:      util.GetGoMod(),
		ControllerName:   values[0],
		UsecaseName:      values[1],
		InportMethodName: "Execute",
	}, nil
}

func (obj *ControllerModel) Run() error {

	InitiateError()
	InitiateLog()
	InitiateHelper()

	err := util.CreateFolderIfNotExist("controller/%s", strings.ToLower(obj.ControllerName))
	if err != nil {
		return err
	}

	{
		fileReadPath := fmt.Sprintf("usecase/%s/port/inport.go", obj.UsecaseName)

		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, fileReadPath, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		for _, decl := range astFile.Decls {

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

				if ts.Name.String() != fmt.Sprintf("%sInport", obj.UsecaseName) {
					continue
				}

				for _, meths := range iFace.Methods.List {

					obj.InportMethodName = meths.Names[0].String()
					break

				}
			}
		}
	}

	{
		outputFile := fmt.Sprintf("controller/response.go")
		err = util.WriteFileIfNotExist(templates.ControllerResponseFile, outputFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(outputFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(outputFile, newBytes, 0644); err != nil {
			return err
		}
	}

	{
		outputFile := fmt.Sprintf("controller/%s/%s.go", strings.ToLower(obj.ControllerName), strings.ToLower(obj.UsecaseName))
		err = util.WriteFileIfNotExist(templates.ControllerGinFile, outputFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(outputFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(outputFile, newBytes, 0644); err != nil {
			return err
		}
	}

	return nil

}

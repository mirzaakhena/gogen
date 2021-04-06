package gogencommand

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
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

	err2 := obj.getInportName()
	if err2 != nil {
		return err2
	}

	responseFile := fmt.Sprintf("controller/%s/response.go", strings.ToLower(obj.ControllerName))
	if !util.IsExist(responseFile) {
		err = util.WriteFile(templates.ControllerResponseFile, responseFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(responseFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(responseFile, newBytes, 0644); err != nil {
			return err
		}
	}

	interceptorFile := fmt.Sprintf("controller/%s/interceptor.go", strings.ToLower(obj.ControllerName))
	if !util.IsExist(interceptorFile) {
		err = util.WriteFile(templates.ControllerInterceptorGinFile, interceptorFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(interceptorFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(interceptorFile, newBytes, 0644); err != nil {
			return err
		}
	}


	controllerFile := fmt.Sprintf("controller/%s/controller.go", strings.ToLower(obj.ControllerName))
	if !util.IsExist(controllerFile) {
		err = util.WriteFile(templates.ControllerGinFile, controllerFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(controllerFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(controllerFile, newBytes, 0644); err != nil {
			return err
		}
	}

	controllerHandlerFile := fmt.Sprintf("controller/%s/handler_%s.go", strings.ToLower(obj.ControllerName), strings.ToLower(obj.UsecaseName))
	if !util.IsExist(controllerHandlerFile) {
		err = util.WriteFileIfNotExist(templates.ControllerFuncGinFile, controllerHandlerFile, obj)
		if err != nil {
			return err
		}

		newBytes, err := imports.Process(controllerHandlerFile, nil, nil)
		if err != nil {
			return err
		}

		if err := ioutil.WriteFile(controllerHandlerFile, newBytes, 0644); err != nil {
			return err
		}

		{

			inportLine, err := obj.getInportLine(controllerFile)
			if err != nil {
				return err
			}

			templateCode, err := util.PrintTemplate(templates.ControllerInportFile, obj)
			if err != nil {
				return err
			}

			file, err := os.Open(controllerFile)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(file)
			var buffer bytes.Buffer
			line := 0
			for scanner.Scan() {
				row := scanner.Text()

				if line == inportLine-1 {
					buffer.WriteString(templateCode)
					buffer.WriteString("\n")
				}

				buffer.WriteString(row)
				buffer.WriteString("\n")
				line++
			}
			if err := file.Close(); err != nil {
				return err
			}

			// reformat and do import
			newBytes, err := imports.Process(controllerFile, buffer.Bytes(), nil)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(controllerFile, newBytes, 0644); err != nil {
				return err
			}

		}

		{

			routerLine, err := obj.getBindRouterLine(controllerFile)
			if err != nil {
				return err
			}

			templateCode, err := util.PrintTemplate(templates.ControllerBindRouterGinFile, obj)
			if err != nil {
				return err
			}

			file, err := os.Open(controllerFile)
			if err != nil {
				return err
			}

			scanner := bufio.NewScanner(file)
			var buffer bytes.Buffer
			line := 0
			for scanner.Scan() {
				row := scanner.Text()

				if line == routerLine-1 {
					buffer.WriteString(templateCode)
					buffer.WriteString("\n")
				}

				buffer.WriteString(row)
				buffer.WriteString("\n")
				line++
			}
			if err := file.Close(); err != nil {
				return err
			}

			// reformat and do import
			newBytes, err := imports.Process(controllerFile, buffer.Bytes(), nil)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(controllerFile, newBytes, 0644); err != nil {
				return err
			}


		}

	}

	return nil

}

func (obj *ControllerModel) getBindRouterLine(controllerFile string) (int, error) {
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}
	routerLine := 0
	for _, decl := range astFile.Decls {

		if gen, ok := decl.(*ast.FuncDecl); ok {

			if gen.Recv == nil {
				continue
			}

			startExp, ok := gen.Recv.List[0].Type.(*ast.StarExpr)
			if !ok {
				continue
			}

			if startExp.X.(*ast.Ident).String() != "Controller" {
				continue
			}

			if gen.Name.String() != "RegisterRouter" {
				continue
			}

			routerLine = fset.Position(gen.Body.Rbrace).Line
			return routerLine, nil
		}

	}
	return 0, fmt.Errorf("register router Not found")
}

func (obj *ControllerModel) getInportLine(controllerFile string) (int, error) {

	inportLine := 0
	fset := token.NewFileSet()
	astFile, err := parser.ParseFile(fset, controllerFile, nil, parser.ParseComments)
	if err != nil {
		return 0, err
	}

	// loop the outport for imports
	for _, decl := range astFile.Decls {

		if gen, ok := decl.(*ast.GenDecl); ok {

			if gen.Tok != token.TYPE {
				continue
			}

			for _, specs := range gen.Specs {

				ts, ok := specs.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if iStruct, ok := ts.Type.(*ast.StructType); ok {

					// check the specific struct name
					if ts.Name.String() != "Controller" {
						continue
					}

					inportLine = fset.Position(iStruct.Fields.Closing).Line
					return inportLine, nil
				}

			}

		}

	}

	return 0, fmt.Errorf(" Controller struct not found")

}

func (obj *ControllerModel) getInportName() error {
	{
		fileReadPath := fmt.Sprintf("usecase/%s/inport.go", obj.UsecaseName)

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

				if ts.Name.String() != "Inport" {
					continue
				}

				for _, meths := range iFace.Methods.List {

					obj.InportMethodName = meths.Names[0].String()
					break

				}
			}
		}
	}
	return nil
}

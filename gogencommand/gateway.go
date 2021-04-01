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

type GatewayModel struct {
	PackagePath string            //
	UsecaseName string            //
	GatewayName string            //
	Methods     []*method         //
	ImportPath  map[string]string //
}

type method struct {
	MethodSignature  string
	DefaultReturnVal string //
}

func NewGatewayModel() (Commander, error) {
	flag.Parse()

	// create a gateway need 2 parameter
	values := flag.Args()[1:]
	if len(values) != 2 {
		return nil, fmt.Errorf("gateway name must not empty. `gogen controller GatewayName UsecaseName`")
	}

	return &GatewayModel{
		PackagePath: util.GetGoMod(),
		GatewayName: values[0],
		UsecaseName: values[1],
		Methods:     []*method{},
		ImportPath:  map[string]string{},
	}, nil
}

func (obj *GatewayModel) Run() error {

	// we need it
	InitiateError()
	InitiateLog()
	InitiateHelper()

	// create a gateway folder
	err := util.CreateFolderIfNotExist("gateway")
	if err != nil {
		return err
	}

	// read outport and collect all the methods
	err = obj.readOutport()
	if err != nil {
		return err
	}

	// write the gateway file
	{

		outputFile := fmt.Sprintf("gateway/%s.go", strings.ToLower(obj.GatewayName))
		err = util.WriteFileIfNotExist(templates.GatewayFile, outputFile, obj)
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

func (obj *GatewayModel) readOutport() error {

	// read the outport.go
	fileReadPath := fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(obj.UsecaseName))

	bytes, err := ioutil.ReadFile(fileReadPath)
	if err != nil {
		return err
	}

	fileStr := string(bytes)

	// split all the line of code into []string
	lineSplit := strings.Split(fileStr, "\n")

	fset := token.NewFileSet()
	astOutportFile, err := parser.ParseFile(fset, fileReadPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// loop the outport for imports
	for _, decl := range astOutportFile.Decls {

		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if gen.Tok == token.IMPORT {
			obj.handleImports(gen, lineSplit, fset)

		}

	}

	// loop the outport for interfaces
	for _, decl := range astOutportFile.Decls {

		gen, ok := decl.(*ast.GenDecl)
		if !ok {
			continue
		}

		if gen.Tok == token.TYPE {
			err = obj.handleInterfaces(gen, lineSplit, fset)
			if err != nil {
				return err
			}

		}

	}

	return nil

}

func (obj *GatewayModel) handleImports(gen *ast.GenDecl, lineSplit []string, fset *token.FileSet) {

	for _, specs := range gen.Specs {

		is, ok := specs.(*ast.ImportSpec)
		if !ok {
			continue
		}

		v := strings.Trim(is.Path.Value, "\"")
		if is.Name != nil {
			obj.ImportPath[is.Name.String()] = v
		} else {
			obj.ImportPath[v[strings.LastIndex(v, "/")+1:]] = v
		}

	}
}

func (obj *GatewayModel) handleInterfaces(gen *ast.GenDecl, lineSplit []string, fset *token.FileSet) error {
	for _, specs := range gen.Specs {

		ts, ok := specs.(*ast.TypeSpec)
		if !ok {
			continue
		}

		iFace, ok := ts.Type.(*ast.InterfaceType)
		if !ok {
			continue
		}

		// check the specific outport interface
		if ts.Name.String() != fmt.Sprintf("%sOutport", obj.UsecaseName) {
			continue
		}

		for _, meths := range iFace.Methods.List {

			if fType, ok := meths.Type.(*ast.SelectorExpr); ok {

				expression := fType.X.(*ast.Ident).String()
				pathWithGomod := obj.ImportPath[expression]
				pathOnly := strings.TrimPrefix(pathWithGomod, obj.PackagePath+"/")

				err := obj.readInterface(fType.Sel.String(), pathOnly)
				if err != nil {
					return err
				}

			} else //

			if fType, ok := meths.Type.(*ast.FuncType); ok {

				// if this is direct method in interface, then handle it
				ms := lineSplit[fset.Position(meths.Names[0].NamePos).Line-1]
				err := obj.handleDefaultReturnValues(fType, ms)
				if err != nil {
					return err
				}

			}

		}

	}

	return nil
}

func (obj *GatewayModel) readInterface(interfaceName, fileReadPath string) error {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, fileReadPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {

			// read the repository.go
			bytes, err := ioutil.ReadFile(fset.File(file.Package).Name())
			if err != nil {
				return err
			}

			fileStr := string(bytes)

			lineSplit := strings.Split(fileStr, "\n")

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

					if ts.Name.String() != interfaceName {
						continue
					}

					for _, meths := range iFace.Methods.List {

						// currently only expect the function
						if fType, ok := meths.Type.(*ast.FuncType); ok {

							ms := lineSplit[fset.Position(meths.Names[0].NamePos).Line-1]
							err = obj.handleDefaultReturnValues(fType, ms)
							if err != nil {
								return err
							}
						}

					}

				}

			}
		}
	}

	return nil

}

func (obj *GatewayModel) handleDefaultReturnValues(fType *ast.FuncType, methodSignature string) error {

	ms := strings.TrimSpace(methodSignature)

	errMsg := fmt.Errorf("function `%s` must have context.Context in its first param argument", ms)

	if fType.Params == nil {
		return errMsg
	}

	if len(fType.Params.List) == 0 {
		return errMsg
	}

	se, ok := fType.Params.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return errMsg
	}

	if fmt.Sprintf("%s.%s", se.X.(*ast.Ident).String(), se.Sel.String()) != "context.Context" {
		return errMsg
	}

	defRetVal := ""
	lenRetList := len(fType.Results.List)
	for i, retList := range fType.Results.List {

		v := ""
		switch t := retList.Type.(type) {

		case *ast.SelectorExpr:
			v = fmt.Sprintf("%v.%v{}", t.X, t.Sel)

		case *ast.StarExpr:
			v = "nil"

		case *ast.Ident:
			if t.String() == "error" {
				v = "nil"
			} else //

			if strings.HasPrefix(t.String(), "int") {
				v = "0"
			} else //

			if t.String() == "string" {
				v = "\"\""
			} else //

			if strings.HasPrefix(t.String(), "float") {
				v = "0.0"
			} else //

			if t.String() == "bool" {
				v = "false"

			} else //

			{
				v = "nil"
			}

		default:
			v = "nil"

		}

		// append the comma
		if i < lenRetList-1 {
			defRetVal += v + ", "
		} else {
			defRetVal += v
		}

	}

	obj.Methods = append(obj.Methods, &method{
		MethodSignature:  ms,
		DefaultReturnVal: defRetVal,
	})

	return nil
}

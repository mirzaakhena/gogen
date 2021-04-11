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

type GatewayModel struct {
	PackagePath       string            //
	UsecaseName       string            //
	GatewayName       string            //
	Methods           []*method         //
	ImportOutportPath map[string]string //
}

type method struct {
	MethodName       string //
	MethodSignature  string //
	DefaultReturnVal string //
}

func NewGatewayModel() (Commander, error) {
	flag.Parse()

	// create a gateway need 2 parameter
	values := flag.Args()[1:]
	if len(values) != 2 {
		return nil, fmt.Errorf("gateway name must not empty. `gogen gateway GatewayName UsecaseName`")
	}

	return &GatewayModel{
		PackagePath:       util.GetGoMod(),
		GatewayName:       values[0],
		UsecaseName:       values[1],
		Methods:           []*method{},
		ImportOutportPath: map[string]string{},
	}, nil
}

func (obj *GatewayModel) Run() error {

	err := InitiateError()
	if err != nil {
		return err
	}

	err = InitiateLog()
	if err != nil {
		return err
	}

	err = InitiateHelper()
	if err != nil {
		return err
	}

	// create a gateway folder
	err = util.CreateFolderIfNotExist("gateway")
	if err != nil {
		return err
	}

	// read outport and collect all the methods
	// we want to create gateway
	// so we start by reading the usecase's outport
	err = obj.readOutport()
	if err != nil {
		return err
	}

	gatewayFile := fmt.Sprintf("gateway/%s.go", strings.ToLower(obj.GatewayName))

	// gateway file not exist yet
	if !util.IsExist(gatewayFile) {
		err := util.WriteFile(templates.GatewayFile, gatewayFile, obj)
		if err != nil {
			return err
		}

		// write the gateway file
		{

			newBytes, err := imports.Process(gatewayFile, nil, nil)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(gatewayFile, newBytes, 0644); err != nil {
				return err
			}

		}

	} else //

	// already exist. check existency current function implementation
	{

		// we will read the existing gateway file
		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, gatewayFile, nil, parser.ParseComments)
		if err != nil {
			return err
		}

		//ast.Print(fset, astFile)

		// prepare the existing function container
		existingFunc := map[string]int{}
		for _, decl := range astFile.Decls {

			// first we want to collect the import only
			gd, ok := decl.(*ast.GenDecl)
			if !ok || gd.Tok.String() != "import" {
				continue
			}

			importPathGateway := map[string]string{}
			handleImports(gd, importPathGateway)

			//fmt.Printf("%v\n", importPathGateway)

			for _, spec := range gd.Specs {
				ts, ok := spec.(*ast.TypeSpec)
				if !ok {
					continue
				}

				if ts.Name.String() != fmt.Sprintf("%sGateway", util.CamelCase(obj.GatewayName)) {
					continue
				}

				st, ok := ts.Type.(*ast.StructType)
				if !ok {
					continue
				}

				if st.Fields.List == nil {
					break
				}

				for _, fieldList := range st.Fields.List {
					se, ok := fieldList.Type.(*ast.SelectorExpr)
					if !ok {
						continue
					}

					pathWithGomod := importPathGateway[se.X.(*ast.Ident).String()]
					pathOnly := strings.TrimPrefix(pathWithGomod, obj.PackagePath+"/")
					structName := se.Sel.String()
					err := obj.readStruct(structName, pathOnly)

					if err != nil {
						return err
					}

				}

			} //

			if !obj.findAndCollectImplMethod(decl, existingFunc) {
				continue
			}
		}

		// collect the only methods that has not added yet
		notExistingMethod := []*method{}
		for _, m := range obj.Methods {
			if _, exist := existingFunc[m.MethodName]; !exist {
				notExistingMethod = append(notExistingMethod, m)
			}
		}

		// we will only inject the non existing method
		obj.Methods = notExistingMethod
		{

			// reopen the file
			file, err := os.Open(gatewayFile)
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

			constTemplateCode, err := util.PrintTemplate(templates.GatewayMethodFile, obj)
			if err != nil {
				return err
			}

			// write the template in the end of file
			buffer.WriteString(constTemplateCode)
			buffer.WriteString("\n")

			// reformat and do import
			newBytes, err := imports.Process(gatewayFile, buffer.Bytes(), nil)
			if err != nil {
				return err
			}

			if err := ioutil.WriteFile(gatewayFile, newBytes, 0644); err != nil {
				return err
			}

		}

	}

	return nil

}

func (obj *GatewayModel) findAndCollectImplMethod(decl ast.Decl, existingFunc map[string]int) bool {
	fd, ok := decl.(*ast.FuncDecl)
	if !ok || fd.Recv == nil {
		return false
	}

	// read all the function that have receiver with gateway name
	if fd.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).String() != fmt.Sprintf("%sGateway", util.CamelCase(obj.GatewayName)) {
		return false
	}

	// collect all existing function that have been there in the file
	existingFunc[fd.Name.String()] = 1

	return true
}

func (obj *GatewayModel) readOutport() error {

	// the Outport interface is under the specific usecase
	// the filename is not restricted.
	// so we will scan all the file under specific usecase folder
	fileReadPath := fmt.Sprintf("usecase/%s", strings.ToLower(obj.UsecaseName))
	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, fileReadPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {

		// read file by file
		for _, file := range pkg.Files {

			//ast.Print(fset, file)

			// for each file, we will read line by line
			for _, decl := range file.Decls {

				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				// we (re)initialize import path for the file
				importPath := map[string]string{}
				handleImports(gen, importPath)

				for _, spec := range gen.Specs {

					// Outport is a type spec
					ts, ok := spec.(*ast.TypeSpec)

					// start by looking the Outport interface with name "Outport"
					if !ok || ts.Name.String() != "Outport" {
						continue
					}

					// make sure Outport is an interface
					iFace, ok := ts.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}

					ast.Print(fset, iFace)

					// start by look up the field
					for _, field := range iFace.Methods.List {

						// as a field, there are two possibility
						switch ty := field.Type.(type) {

						case *ast.SelectorExpr: // as interface extension
							fmt.Print("%v", ty)


						case *ast.FuncType: // as direct func (method) interface
							fmt.Print("as a function %v", ty)

						}

					}

				}

			}

			//// loop the outport for imports
			//for _, decl := range file.Decls {
			//
			//	gen, ok := decl.(*ast.GenDecl)
			//	if !ok {
			//		continue
			//	}
			//
			//	if gen.Tok == token.IMPORT {
			//		handleImports(gen, obj.ImportOutportPath)
			//
			//	}
			//
			//}
			//
			//port := file.Name.String()
			//
			//// loop the outport for interfaces
			//for _, decl := range file.Decls {
			//
			//	gen, ok := decl.(*ast.GenDecl)
			//	if !ok {
			//		continue
			//	}
			//
			//	if gen.Tok == token.TYPE {
			//		err = obj.handleInterfaces(port, gen)
			//		if err != nil {
			//			return err
			//		}
			//
			//	}
			//
			//}

		}
	}

	// read the outport.go

	//fset := token.NewFileSet()
	//astOutportFile, err := parser.ParseFile(fset, fileReadPath, nil, parser.ParseComments)
	//if err != nil {
	//	return err
	//}

	return nil

}

func handleImports(gen *ast.GenDecl, ImportPath map[string]string) {

	for _, specs := range gen.Specs {

		is, ok := specs.(*ast.ImportSpec)
		if !ok {
			continue
		}

		v := strings.Trim(is.Path.Value, "\"")
		if is.Name != nil {
			ImportPath[is.Name.String()] = v
		} else {
			ImportPath[v[strings.LastIndex(v, "/")+1:]] = v
		}

	}
}

func (obj *GatewayModel) handleInterfaces(port string, gen *ast.GenDecl) error {
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
		if ts.Name.String() != "Outport" {
			continue
		}

		for _, meths := range iFace.Methods.List {

			// extend another interface
			if selectorExp, ok := meths.Type.(*ast.SelectorExpr); ok {

				expression := selectorExp.X.(*ast.Ident).String()
				pathWithGomod := obj.ImportOutportPath[expression]
				pathOnly := strings.TrimPrefix(pathWithGomod, obj.PackagePath+"/")
				err := obj.readInterface(selectorExp.Sel.String(), pathOnly)
				if err != nil {
					return err
				}

			} else
			// direct function in interface
			if fType, ok := meths.Type.(*ast.FuncType); ok {

				// if this is direct method in interface, then handle it
				err := obj.handleMethodSignature(port, fType, meths.Names[0].String())
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

			port := file.Name.String()

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

							err = obj.handleMethodSignature(port, fType, meths.Names[0].String())
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

func (obj *GatewayModel) readStruct(structName, fileReadPath string) error {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, fileReadPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {

			ast.Print(fset, file)

			//port := file.Name.String()
			//
			//for _, decl := range file.Decls {
			//
			//	gen, ok := decl.(*ast.GenDecl)
			//	if !ok {
			//		continue
			//	}
			//
			//	if gen.Tok != token.TYPE {
			//		continue
			//	}
			//
			//	for _, specs := range gen.Specs {
			//
			//		ts, ok := specs.(*ast.TypeSpec)
			//		if !ok {
			//			continue
			//		}
			//
			//		iFace, ok := ts.Type.(*ast.StructType)
			//		if !ok {
			//			continue
			//		}
			//
			//		if ts.Name.String() != structName {
			//			continue
			//		}
			//
			//		for _, meths := range iFace.Fields.List {
			//
			//			// currently only expect the function
			//			if fType, ok := meths.Type.(*ast.FuncType); ok {
			//
			//				err = obj.handleMethodSignature(port, fType, meths.Names[0].String())
			//				if err != nil {
			//					return err
			//				}
			//			}
			//
			//		}
			//
			//	}
			//
			//}

		}
	}

	return nil

}

func (obj *GatewayModel) handleMethodSignature(port string, fType *ast.FuncType, methodName string) error {

	ms := strings.TrimSpace(methodName)

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
	if fType.Results != nil {
		lenRetList := len(fType.Results.List)
		for i, retList := range fType.Results.List {

			v := ""
			switch t := retList.Type.(type) {

			case *ast.SelectorExpr:
				v = fmt.Sprintf("%v.%v{}", t.X, t.Sel)

			case *ast.StarExpr:
				v = "nil"

			case *ast.Ident:

				if t.Name == "error" {
					v = "nil"
				} else //

				if strings.HasPrefix(t.Name, "int") {
					v = "0"
				} else //

				if t.Name == "string" {
					v = "\"\""
				} else //

				if strings.HasPrefix(t.Name, "float") {
					v = "0.0"
				} else //

				if t.Name == "bool" {
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

	}

	methodSignature := FuncHandler{Port: port}.processFuncType(&bytes.Buffer{}, fType)
	obj.Methods = append(obj.Methods, &method{
		MethodName:       methodName,
		MethodSignature:  methodSignature,
		DefaultReturnVal: defRetVal,
	})

	return nil
}

//===============================

type FuncHandler struct {
	Port string
}

func (r FuncHandler) appendType(expr ast.Expr) string {
	var param bytes.Buffer

	for {
		switch t := expr.(type) {
		case *ast.Ident:
			return r.processIdent(&param, t)

		case *ast.ArrayType:
			return r.processArrayType(&param, t)

		case *ast.StarExpr:
			return r.processStarExpr(&param, t)

		case *ast.SelectorExpr:
			return r.processSelectorExpr(&param, t)

		case *ast.InterfaceType:
			return r.processInterfaceType(&param, t)

		case *ast.ChanType:
			return r.processChanType(&param, t)

		case *ast.MapType:
			return r.processMapType(&param, t)

		case *ast.FuncType:
			param.WriteString("func")
			return r.processFuncType(&param, t)

		case *ast.StructType:
			return r.processStruct(&param, t)

		default:
			return param.String()
		}

	}
}

func (r FuncHandler) processIdent(param *bytes.Buffer, t *ast.Ident) string {
	if t.Obj != nil {
		param.WriteString(r.Port)
		param.WriteString(".")
	}
	param.WriteString(t.Name)
	return param.String()
}

func (r FuncHandler) processStruct(param *bytes.Buffer, t *ast.StructType) string {
	param.WriteString("struct{")

	nParam := t.Fields.NumFields()
	c := 0
	for _, field := range t.Fields.List {
		nNames := len(field.Names)
		c += nNames
		for i, name := range field.Names {
			param.WriteString(name.String())
			if i < len(field.Names)-1 {
				param.WriteString(", ")
			} else {
				param.WriteString(" ")
			}
		}
		param.WriteString(r.appendType(field.Type))
		//fmt.Printf(">>> %v:%v:%v:%v\n", nParam, iList, nNames, c)

		if c < nParam {
			param.WriteString("; ")
		}

	}

	param.WriteString("}")
	return param.String()
}

func (r FuncHandler) processArrayType(param *bytes.Buffer, t *ast.ArrayType) string {
	if t.Len != nil {
		arrayCapacity := t.Len.(*ast.BasicLit).Value
		param.WriteString(fmt.Sprintf("[%s]", arrayCapacity))
	} else {
		param.WriteString("[]")
	}
	param.WriteString(r.appendType(t.Elt))

	return param.String()
}

func (r FuncHandler) processStarExpr(param *bytes.Buffer, t *ast.StarExpr) string {
	param.WriteString("*")
	param.WriteString(r.appendType(t.X))
	return param.String()
}

func (r FuncHandler) processInterfaceType(param *bytes.Buffer, t *ast.InterfaceType) string {
	param.WriteString("interface{}")
	_ = t
	return param.String()
}

func (r FuncHandler) processSelectorExpr(param *bytes.Buffer, t *ast.SelectorExpr) string {
	param.WriteString(r.appendType(t.X))
	param.WriteString(".")
	param.WriteString(t.Sel.Name)
	return param.String()
}

func (r FuncHandler) processChanType(param *bytes.Buffer, t *ast.ChanType) string {
	if t.Dir == 1 {
		param.WriteString("chan<- ")
	} else if t.Dir == 2 {
		param.WriteString("<-chan ")
	} else {
		param.WriteString("chan ")
	}

	param.WriteString(r.appendType(t.Value))
	return param.String()
}

func (r FuncHandler) processMapType(param *bytes.Buffer, t *ast.MapType) string {
	param.WriteString("map[")
	param.WriteString(r.appendType(t.Key))
	param.WriteString("]")
	param.WriteString(r.appendType(t.Value))
	return param.String()
}

func (r FuncHandler) processFuncType(param *bytes.Buffer, t *ast.FuncType) string {

	nParam := t.Params.NumFields()
	param.WriteString("(")
	for iList, field := range t.Params.List {
		nNames := len(field.Names)
		for i, name := range field.Names {
			param.WriteString(name.String())
			if i < len(field.Names)-1 {
				param.WriteString(", ")
			} else {
				param.WriteString(" ")
			}
		}
		param.WriteString(r.appendType(field.Type))
		if iList+nNames < nParam {
			param.WriteString(", ")
		}

	}
	param.WriteString(") ")

	nResult := t.Results.NumFields()

	haveParentThesis := false
	if t.Results == nil {
		return param.String()
	}
	for i, field := range t.Results.List {

		if i == 0 {
			if len(field.Names) > 0 || nResult > 1 {
				param.WriteString("(")
				haveParentThesis = true
			}
		}

		for i, name := range field.Names {
			param.WriteString(name.String() + " ")
			if i < len(field.Names)-1 {
				param.WriteString(", ")
			}
		}
		param.WriteString(r.appendType(field.Type))
		if i+1 < nResult {
			param.WriteString(", ")
		}
	}

	if haveParentThesis {
		param.WriteString(")")
	}

	return param.String()
}

func (r *GatewayModel) readFunction() {

}

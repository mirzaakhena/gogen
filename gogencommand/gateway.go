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
	PackagePath string    //
	UsecaseName string    //
	GatewayName string    //
	Methods     []*method // all the method needed
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
		PackagePath: util.GetGoMod(),
		GatewayName: values[0],
		UsecaseName: values[1],
		Methods:     []*method{},
		//ImportOutportPath: map[string]string{},
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

	err = obj.ReadOutport()
	if err != nil {
		return err
	}

	// create a gateway folder
	err = util.CreateFolderIfNotExist("gateway")
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

		existingFunc, err := obj.readCurrentGateway()
		if err != nil {
			return err
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

//===============================

func (obj *GatewayModel) ReadOutport() error {
	fileReadPath := fmt.Sprintf("usecase/%s", strings.ToLower(obj.UsecaseName))

	err := obj.readInterface("Outport", fileReadPath)
	if err != nil {
		return err
	}

	return nil
}

func (obj *GatewayModel) readInterface(interfaceName, folderPath string) error {

	packagePath := util.GetGoMod()

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {

		// read file by file
		for _, file := range pkg.Files {

			ip := map[string]string{}

			// for each file, we will read line by line
			for _, decl := range file.Decls {

				gen, ok := decl.(*ast.GenDecl)
				if !ok {
					continue
				}

				for _, spec := range gen.Specs {

					is, ok := spec.(*ast.ImportSpec)
					if ok {
						handleImports(is, ip)
					}

					// Outport is must a type spec
					ts, ok := spec.(*ast.TypeSpec)
					if !ok {
						continue
					}

					// start by looking the Outport interface with name "Outport"
					if ts.Name.String() != interfaceName {
						continue
					}

					// make sure Outport is an interface
					iFace, ok := ts.Type.(*ast.InterfaceType)
					if !ok {
						continue
					}

					for _, field := range iFace.Methods.List {

						// as a field, there are two possibility
						switch ty := field.Type.(type) {

						case *ast.SelectorExpr: // as extension of another interface
							expression := ty.X.(*ast.Ident).String()
							pathWithGomod := ip[expression]
							pathOnly := strings.TrimPrefix(pathWithGomod, packagePath+"/")
							interfaceName := ty.Sel.String()
							err := obj.readInterface(interfaceName, pathOnly)
							if err != nil {
								return err
							}

						case *ast.FuncType: // as direct func (method) interface
							//TODO cannot handle Something(c context.Context, a Hoho) yet, the Hoho part
							msObj, err := obj.handleMethodSignature(file.Name.String(), ty, field.Names[0].String())
							if err != nil {
								return err
							}
							obj.Methods = append(obj.Methods, msObj)

						case *ast.Ident: // as interface extension in same package
							//ast.Print(fset, ty)
							//TODO as interface extension in same package in the same or different file
							fmt.Printf("as extension in same import %v\n", ty)
						}

					}

				}

			}

		}
	}

	return nil
}

func (obj *GatewayModel) handleMethodSignature(prefixExpression string, fType *ast.FuncType, methodName string) (*method, error) {

	ms := strings.TrimSpace(methodName)

	errMsg := fmt.Errorf("function `%s` must have context.Context in its first param argument", ms)

	if fType.Params == nil || len(fType.Params.List) == 0 {
		return nil, errMsg
	}

	se, ok := fType.Params.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return nil, errMsg
	}

	if fmt.Sprintf("%s.%s", se.X.(*ast.Ident).String(), se.Sel.String()) != "context.Context" {
		return nil, errMsg
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

				} else if strings.HasPrefix(t.Name, "int") {
					v = "0"

				} else if t.Name == "string" {
					v = "\"\""

				} else if strings.HasPrefix(t.Name, "float") {
					v = "0.0"

				} else if t.Name == "bool" {
					v = "false"

				} else {
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

	methodSignature := FuncHandler{PrefixExpression: prefixExpression}.processFuncType(&bytes.Buffer{}, fType)
	msObj := method{
		MethodName:       methodName,
		MethodSignature:  methodSignature,
		DefaultReturnVal: defRetVal,
	}

	return &msObj, nil
}

func handleImports(is *ast.ImportSpec, ip map[string]string) {
	v := strings.Trim(is.Path.Value, "\"")
	if is.Name != nil {
		ip[is.Name.String()] = v
	} else {
		ip[v[strings.LastIndex(v, "/")+1:]] = v
	}
}

// --------------------------------------

func (obj *GatewayModel) readCurrentGateway() (map[string]int, error) {

	structName := fmt.Sprintf("%sGateway", util.CamelCase(obj.GatewayName))
	fileReadPath := fmt.Sprintf("gateway/")

	existingFunc := map[string]int{}

	err := obj.readStruct(structName, fileReadPath, existingFunc)
	if err != nil {
		return nil, err
	}

	return existingFunc, nil
}

func (obj *GatewayModel) readStruct(structName, folderPath string, existingFunc map[string]int) error {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, folderPath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {

		// read file by file
		for _, file := range pkg.Files {

			importPaths := map[string]string{}

			for _, decl := range file.Decls {

				switch gd := decl.(type) {

				case *ast.GenDecl:
					err := obj.generalDecl(structName, gd, importPaths, existingFunc)
					if err != nil {
						return err
					}

				case *ast.FuncDecl:
					//ast.Print(fset, gd)
					if !obj.findAndCollectImplMethod(gd, structName, existingFunc) {
						continue
					}
				}

			}

		}

	}

	return nil
}

func (obj *GatewayModel) generalDecl(structName string, gd *ast.GenDecl, importPaths map[string]string, existingFunc map[string]int) error {
	for _, spec := range gd.Specs {

		// handle import
		is, ok := spec.(*ast.ImportSpec)
		if ok {
			handleImports(is, importPaths)
		}

		// it is type declaration
		ts, ok := spec.(*ast.TypeSpec)
		if !ok {
			continue
		}

		// the struct name must have a 'Gateway' suffix
		if ts.Name.String() != structName {
			continue
		}

		// gateway must be a struct type
		st, ok := ts.Type.(*ast.StructType)
		if !ok {
			continue
		}

		// if struct list empty then nothing to do
		if st.Fields.List == nil {
			break
		}

		for _, fieldList := range st.Fields.List {

			switch ty := fieldList.Type.(type) {
			case *ast.SelectorExpr: // struct is extend another struct

				expression := ty.X.(*ast.Ident).String()
				pathWithGomod := importPaths[expression]
				pathOnly := strings.TrimPrefix(pathWithGomod, obj.PackagePath+"/")
				structName := ty.Sel.String()
				err := obj.readStruct(structName, pathOnly, existingFunc)
				if err != nil {
					return err
				}

			}

		}

	}
	return nil
}

func (obj *GatewayModel) findAndCollectImplMethod(fd *ast.FuncDecl, structName string, existingFunc map[string]int) bool {
	if fd.Recv == nil {
		return false
	}

	// read all the function that have receiver with gateway name
	if fd.Recv.List[0].Type.(*ast.StarExpr).X.(*ast.Ident).String() != structName {
		return false
	}

	// collect all existing function that have been there in the file
	existingFunc[fd.Name.String()] = 1

	return true
}

//===============================

type FuncHandler struct {
	PrefixExpression string //
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
		param.WriteString(r.PrefixExpression)
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

	// TODO need to handle method param without variable
	// TODO need to handle param with struct/interface type

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
		//param.WriteString(r.PrefixExpression)
		//param.WriteString(".")
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

package gengateway2

import (
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"strings"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	DomainName  string
	GatewayName string
	UsecaseName *string
	Methods     utils.OutportMethods
}

const defaultOutportName = "Outport"

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a gateway for all usecases with cloverdb sample implementation\n" +
			"   gogen gateway inmemory\n" +
			"     'inmemory' is a gateway name\n" +
			"\n" +
			"   # Create a gateway for specific usecase\n" +
			"   gogen gateway inmemory cloverdb\n" +
			"     'inmemory' is a gateway name\n" +
			"     'cloverdb' is a sample implementation\n" +
			"\n" +
			"   # Create a gateway for specific usecase\n" +
			"   gogen gateway inmemory cloverdb CreateOrder\n" +
			"     'inmemory'    is a gateway name\n" +
			"     'cloverdb' is a sample implementation\n" +
			"     'CreateOrder' is an usecase name\n" +
			"\n")

		return err
	}

	domainName := utils.GetDefaultDomain()
	gatewayName := inputs[0]

	obj := ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		GatewayName: gatewayName,
		DomainName:  utils.LowerCase(domainName),
		UsecaseName: nil,
	}

	//driverName := "simple"

	//if len(inputs) >= 2 {
	//	driverName = inputs[1]
	//}

	if len(inputs) >= 3 {
		obj.UsecaseName = &inputs[2]
	}

	fileInfo, err := os.ReadDir(fmt.Sprintf("domain_%s/usecase", domainName))
	if err != nil {
		return err
	}

	gatewayObject := GatewayObject{
		UsecaseName: "",
		Imports:     map[ImportAlias]ImportPath{},
		Interfaces:  map[ImportAlias]InterfaceName{},
		Methods:     map[ImportAlias]Method{},
	}

	for _, file := range fileInfo {

		// skip all the file
		if !file.IsDir() {
			continue
		}

		// folder name is a usecase name
		folderName := file.Name()

		//fmt.Printf(">>>> %v\n", folderName)

		packagePath := fmt.Sprintf("domain_%s/usecase/%s", domainName, folderName)

		err := readOutportInterface(packagePath, defaultOutportName, &gatewayObject)
		if err != nil {
			return err
		}

	}

	fmt.Printf(">>> %v\n", gatewayObject)

	return nil

}

func readOutportInterface(packagePath, outportName string, gatewayObject *GatewayObject) error {

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, packagePath, nil, parser.ParseComments)
	if err != nil {
		return err
	}

	// from all package
	for _, pkg := range pkgs {

		foundOutport := false

		// from all files
		for _, file := range pkg.Files {

			fmt.Printf(">>>> %v\n", file.Name.String())

			//ast.Print(fset, file)

			// prepare the temporary import map
			importMap := map[ImportAlias]ImportPath{}

			// walk through it
			ast.Inspect(file, func(node ast.Node) bool {

				// outport already found nothing to do anymore so we can skip it
				if foundOutport {
					return false
				}

				// read all the import
				genDecl, ok := node.(*ast.GenDecl)
				if ok && genDecl.Tok == token.IMPORT {

					for _, spec := range genDecl.Specs {
						importSpec, ok := spec.(*ast.ImportSpec)
						if !ok {
							break
						}

						// from "yourpath/domain_core/model/repository"
						// we want just yourpath/domain_core/model/repository
						// without "
						value := strings.Trim(importSpec.Path.Value, "\"")

						// prefix can be exist or nil
						// exist     : myrepo yourpath/domain_core/model/repository
						// not exist : yourpath/domain_core/model/repository
						prefix := importSpec.Name

						// prefix exist is something like
						// alias "yourpath/domain_core/model/repository"
						// the prefix is "alias"
						if prefix != nil {
							importMap[ImportAlias(prefix.String())] = ImportPath(value)
							continue
						}

						// the prefix is not exist then we use the last path as prefix
						// from "yourpath/domain_core/model/repository"
						// we will get "repository" as the prefix
						// TODO we are not covered the suffix with version yet
						lastIndex := strings.LastIndex(value, "/")
						if lastIndex == -1 {
							// TODO something goes wrong here
							break
						}

						// from yourpath/domain_core/model/repository
						// we use "repository" as the prefix
						newPrefix := value[lastIndex+1:]
						importMap[ImportAlias(newPrefix)] = ImportPath(value)

					}

					return true
				}

				for {

					// find a type
					typeSpec, ok := node.(*ast.TypeSpec)
					if !ok || typeSpec.Name.String() != outportName {
						break
					}

					// type interface
					interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
					if !ok {
						break
					}

					// look up all the method in the interface
					for _, method := range interfaceType.Methods.List {

						switch methodType := method.Type.(type) {

						case *ast.SelectorExpr: // as extension of another interface from external package

							// example : repository.FindAllOrder

							// get the prefix name : repository
							ident, ok := methodType.X.(*ast.Ident)
							if !ok {
								break
							}

							// get the interface name : FindAllOrder
							interfaceName := methodType.Sel.String()

							// mix it tobe : repository.FindAllOrder
							interfaceWithPrefix := fmt.Sprintf("%s.%s", ident.String(), interfaceName)

							// store it
							gatewayObject.Interfaces[ImportAlias(ident.String())] = InterfaceName(interfaceWithPrefix)

							// get the package path from package/alias name
							nextPackagePath := importMap[ImportAlias(ident.String())]

							// store the used import only
							_, exist := gatewayObject.Imports[ImportAlias(ident.String())]
							if !exist {
								gatewayObject.Imports[ImportAlias(ident.String())] = nextPackagePath
							}

							// recursively check the interface
							err := readOutportInterface(string(nextPackagePath), interfaceName, gatewayObject)
							if err != nil {
								// TODO something goes wrong here
								return false
							}

						case *ast.FuncType: // as direct func

							// example : FindAllOrder(ctx context.Context, req FindAllOrderRequest) ([]*entity.Order, int64, error)

							err := handleMethodSignature(method.Names[0].String(), methodType, file.Name.String(), gatewayObject)
							if err != nil {
								return false
							}

						case *ast.Ident: // as extension of another interface from internal package
							//fmt.Printf(">> ident\n")
						}

					}

					foundOutport = true

					//fmt.Printf(">>>>> %v %v\n", pkg.Name, typeSpec.Name.String())

					//ast.Print(fset, ident)

					return false
				}

				return true
			})

			//fmt.Printf("%v\n", importMap)

			if foundOutport {
				break
			}

		}

	}

	return nil
}

func handleMethodSignature(methodName string, funcType *ast.FuncType, prefixExpression string, gatewayObject *GatewayObject) error {

	// checking first params as context.Context
	if !validateFirstParamIsContext(funcType) {
		return fmt.Errorf("function `%s` must have context.Context in its first param argument", methodName)
	}

	// result at least have one returned value which is an error
	if funcType.Results == nil {
		return fmt.Errorf("function `%s` result at least have error return value", methodName)
	}

	//defParVal := composeDefaultValue(funcType.Params.List)
	//defRetVal := composeDefaultValue(funcType.Results.List)
	//
	//methodSignature := utils.TypeHandler{PrefixExpression: prefixExpression}.processFuncType(&bytes.Buffer{}, fType)
	//msObj := method{
	//	MethodName:       methodName,
	//	MethodSignature:  methodSignature,
	//	DefaultParamVal:  defParVal,
	//	DefaultReturnVal: defRetVal,
	//}
	//
	//*obj = append(*obj, &msObj)

	return nil
}

func composeDefaultValue(list []*ast.Field) string {

	defRetVal := ""
	for i, retList := range list {

		var v string

		switch t := retList.Type.(type) {

		case *ast.SelectorExpr:
			v = fmt.Sprintf("%v.%v{}", t.X, t.Sel)

			// little hacky for context
			if v == "context.Context{}" {
				v = "ctx"
			}

		case *ast.StarExpr:
			v = "nil"

		case *ast.Ident:

			if t.Name == "error" {
				v = "nil"

			} else if strings.HasPrefix(t.Name, "int") {
				v = "0"

			} else if strings.HasPrefix(t.Name, "float") {
				v = "0.0"

			} else if t.Name == "string" {
				v = "\"\""

			} else if t.Name == "bool" {
				v = "false"

			} else {
				v = "nil"
			}

		default:
			v = "nil"

		}

		// append the comma
		if i < len(list)-1 {
			defRetVal += v + ", "
		} else {
			defRetVal += v
		}

	}
	return defRetVal
}

func validateFirstParamIsContext(funcType *ast.FuncType) bool {

	if funcType.Params == nil || len(funcType.Params.List) == 0 {
		return false
	}

	se, ok := funcType.Params.List[0].Type.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	firstParamArgs := fmt.Sprintf("%s.%s", se.X.(*ast.Ident).String(), se.Sel.String())
	if firstParamArgs != "context.Context" {
		return false
	}

	return true
}

type (
	ImportAlias   string
	ImportPath    string
	InterfaceName string
	Method        struct {
		MethodName       string
		MethodSignature  string
		DefaultParamVal  string
		DefaultReturnVal string
	}
)

type GatewayObject struct {
	UsecaseName string
	Imports     map[ImportAlias]ImportPath
	Methods     map[ImportAlias]Method
	Interfaces  map[ImportAlias]InterfaceName
}

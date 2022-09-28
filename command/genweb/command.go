package genweb

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
	DomainName        string
	UsecaseName       string
	RequestNameTypes  []*NameType
	ResponseNameTypes []*NameType
}

type NameType struct {
	Name string
	Type string
}

func Run(inputs ...string) error {

	//if len(inputs) < 1 {
	//	err := fmt.Errorf("\n" +
	//		"   # Create a web app\n" +
	//		"   gogen web\n" +
	//		"     'Product' is an existing entity name\n" +
	//		"\n")
	//
	//	return err
	//}

	domainName := utils.GetDefaultDomain()

	//entityName := inputs[0]

	controllerName := "restapi"

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/web/", "shared", fileRenamer, struct{}{}, utils.AppTemplates)
	if err != nil {
		return err
	}

	controllerFolderName := fmt.Sprintf("domain_%s/controller/%s", domainName, controllerName)

	fileInfo, err := os.ReadDir(controllerFolderName)
	if err != nil {
		return err
	}

	for _, file := range fileInfo {

		if !strings.HasPrefix(file.Name(), "handler_") {
			continue
		}

		if strings.HasSuffix(file.Name(), ".http") {
			continue
		}

		fset := token.NewFileSet()
		astFile, err := parser.ParseFile(fset, fmt.Sprintf("%s/%s", controllerFolderName, file.Name()), nil, parser.ParseComments)
		if err != nil {
			fmt.Printf("%v\n", err.Error())
			os.Exit(1)
		}

		//ast.Print(fset, astFile)

		for _, decl := range astFile.Decls {

			funcDecl, ok := decl.(*ast.FuncDecl)
			if !ok {
				continue
			}

			methodName := utils.PascalCase(funcDecl.Name.String())
			if !strings.HasSuffix(methodName, "Handler") {
				continue
			}

			i := strings.LastIndex(methodName, "Handler")

			usecaseName := methodName[:i]

			fileRenamer := map[string]string{
				"domainname":  utils.LowerCase(domainName),
				"usecasename": utils.LowerCase(usecaseName),
			}

			obj := &ObjTemplate{
				DomainName:  domainName,
				UsecaseName: usecaseName,
			}

			if strings.HasPrefix(utils.LowerCase(usecaseName), "getall") {

				// TODO read response field

				err := utils.CreateEverythingExactly("templates/web/", "getall", fileRenamer, obj, utils.AppTemplates)
				if err != nil {
					return err
				}

			} else if strings.HasPrefix(utils.LowerCase(usecaseName), "get") {
				err := utils.CreateEverythingExactly("templates/web/", "get", fileRenamer, obj, utils.AppTemplates)
				if err != nil {
					return err
				}

			} else {

				for _, stmt := range funcDecl.Body.List {
					declStmt, ok := stmt.(*ast.DeclStmt)
					if !ok {
						continue
					}
					genDecl, ok := declStmt.Decl.(*ast.GenDecl)
					if !ok {
						continue
					}

					if genDecl.Tok.String() != "type" {
						continue
					}

					for _, spec := range genDecl.Specs {
						typeSpec, ok := spec.(*ast.TypeSpec)
						if !ok {
							continue
						}

						if typeSpec.Name.String() != "request" {
							continue
						}

						structType, ok := typeSpec.Type.(*ast.StructType)
						if !ok {
							continue
						}

						for _, field := range structType.Fields.List {
							varType, ok := field.Type.(*ast.Ident)
							if !ok {
								continue
							}
							for _, name := range field.Names {
								obj.RequestNameTypes = append(obj.RequestNameTypes, &NameType{
									Name: name.String(),
									Type: varType.String(),
								})
							}
						}
					}

				}

				err := utils.CreateEverythingExactly("templates/web/", "run", fileRenamer, obj, utils.AppTemplates)
				if err != nil {
					return err
				}

			}

		}

	}

	return nil

}

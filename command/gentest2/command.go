package gentest

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"

	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	DomainName  string
	UsecaseName string
	//TestName    string
	Methods utils.OutportMethods
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Create a test case file for current usecase\n" +
			"   gogen test normal CreateOrder\n" +
			"     'normal'      is a test case name\n" +
			"     'CreateOrder' is an usecase name\n" +
			"\n")

		return err
	}

	packagePath := utils.GetPackagePath()
	gcfg := utils.GetGogenConfig()
	//testName := inputs[0]
	usecaseName := inputs[0]

	obj := ObjTemplate{
		PackagePath: packagePath,
		UsecaseName: usecaseName,
		//TestName:    testName,
		DomainName: gcfg.Domain,
		Methods:    nil,
	}

	outportMethods, err := utils.NewOutportMethods(gcfg.Domain, usecaseName)
	if err != nil {
		return err
	}

	obj.Methods = outportMethods

	fileRenamer := map[string]string{
		"domainname":  utils.LowerCase(gcfg.Domain),
		"usecasename": utils.LowerCase(usecaseName),
		//"testname":    utils.LowerCase(testName),
	}
	err = utils.CreateEverythingExactly("templates/", "test_new", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	fset := token.NewFileSet()
	pkgs, err := parser.ParseDir(fset, fmt.Sprintf("domain_%s/usecase/%s", utils.LowerCase(gcfg.Domain), utils.LowerCase(usecaseName)), nil, parser.ParseComments)
	if err != nil {
		return err
	}

	for _, pkg := range pkgs {
		for _, file := range pkg.Files {
			ast.Inspect(file, func(node ast.Node) bool {

				// if there is an error, just ignore everything
				if err != nil {
					return false
				}

				typeSpec, ok := node.(*ast.TypeSpec)
				if !ok {
					return true
				}

				typeSpecName := typeSpec.Name.String()

				if typeSpecName != "Outport" {
					return false
				}

				interfaceType, ok := typeSpec.Type.(*ast.InterfaceType)
				if !ok {
					err = fmt.Errorf("outport is not interface")
				}

				for _, method := range interfaceType.Methods.List {
					switch method.Type.(type) {
					case *ast.SelectorExpr:

					}
				}

				return true
			})
		}
	}

	//// file gateway impl file is already exist, we want to inject non existing method
	//structName := fmt.Sprintf("mockOutport%s", utils.PascalCase(testName))
	//rootFolderName := fmt.Sprintf("domain_%s/usecase/%s", utils.LowerCase(gcfg.Domain), utils.LowerCase(usecaseName))
	//existingFunc, err := utils.NewOutportMethodImpl(structName, rootFolderName)
	//if err != nil {
	//	return err
	//}
	//
	//// collect the only methods that has not added yet
	//notExistingMethod := utils.OutportMethods{}
	//for _, m := range outportMethods {
	//	if _, exist := existingFunc[m.MethodName]; !exist {
	//		notExistingMethod = append(notExistingMethod, m)
	//	}
	//}
	//
	//gatewayCode, err := getTestMethodTemplate()
	//if err != nil {
	//	return err
	//}
	//
	//obj.Methods = notExistingMethod
	//
	//templateHasBeenInjected, err := utils.PrintTemplate(string(gatewayCode), obj)
	//if err != nil {
	//	return err
	//}
	//
	//bytes, err := obj.injectToTest(templateHasBeenInjected)
	//if err != nil {
	//	return err
	//}
	//
	//// reformat outport._go
	//err = utils.Reformat(obj.getTestFileName(), bytes)
	//if err != nil {
	//	return err
	//}

	return nil

}

//// getTestMethodTemplate ...
//func getTestMethodTemplate() ([]byte, error) {
//	testfileName := fmt.Sprintf("templates/test/domain_${domainname}/usecase/${usecasename}/~inject._go")
//	return utils.AppTemplates.ReadFile(testfileName)
//}

////func (o ObjTemplate) injectToTest(injectedCode string) ([]byte, error) {
//	return utils.InjectCodeAtTheEndOfFile(o.getTestFileName(), injectedCode)
//}

//// getTestFileName ...
//func (o ObjTemplate) getTestFileName() string {
//	return fmt.Sprintf("domain_%s/usecase/%s/testcase_%s_test.go", utils.LowerCase(o.DomainName), utils.LowerCase(o.UsecaseName), utils.LowerCase(o.TestName))
//}

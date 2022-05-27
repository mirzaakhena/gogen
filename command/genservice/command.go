package genservice

import (
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath string
	ServiceName string
	UsecaseName *string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("   # Create a service and inject the template code into interactor file with '//!' flag\n" +
			"   gogen service PublishMessage CreateOrder\n" +
			"     'PublishMessage' is a service func name\n" +
			"     'CreateOrder'    is an usecase name\n" +
			"\n" +
			"   # Create a service without inject the template code into usecase\n" +
			"   gogen service PublishMessage\n" +
			"     'PublishMessage' is a service func name\n" +
			"\n")
		return err
	}
	domainName := utils.GetDefaultDomain()
	serviceName := inputs[0]

	obj := ObjTemplate{
		PackagePath: utils.GetPackagePath(),
		ServiceName: serviceName,
		UsecaseName: nil,
	}

	if len(inputs) >= 2 {
		obj.UsecaseName = &inputs[1]
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "service", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	// service.go file is already exist, but is the interface is exist ?
	exist, err := isServiceExist(fmt.Sprintf("domain_%s/model/service", domainName), serviceName)
	if err != nil {
		return err
	}

	if !exist {

		// check the prefix and give specific template for it
		templateCode, err := getServiceInterfaceTemplate()
		if err != nil {
			return err
		}
		templateHasBeenInjected, err := utils.PrintTemplate(string(templateCode), obj)
		if err != nil {
			return err
		}

		dataInBytes, err := injectCodeToServiceFile(domainName, templateHasBeenInjected)
		if err != nil {
			return err
		}

		// reformat interactor.go
		err = utils.Reformat(getServiceFile(domainName), dataInBytes)
		if err != nil {
			return err
		}

	}

	// no usecase name means no need to inject to outport and interactor
	if obj.UsecaseName == nil {
		return nil
	}

	// inject to outport
	{
		err := injectToOutport(domainName, obj)
		if err != nil {
			return err
		}

		outportFile := fmt.Sprintf("domain_%s/usecase/%s/outport.go", domainName, *obj.UsecaseName)

		// reformat outport.go
		err = utils.Reformat(outportFile, nil)
		if err != nil {
			return err
		}
	}

	// inject to interactor
	{
		// check the prefix and give specific template for it
		interactorCode, err := getServiceInjectTemplate()
		if err != nil {
			return err
		}

		templateHasBeenInjected, err := utils.PrintTemplate(string(interactorCode), obj)
		if err != nil {
			return err
		}

		interactorFilename := getInteractorFilename(domainName, *obj.UsecaseName)

		interactorBytes, err := utils.InjectToInteractor(interactorFilename, templateHasBeenInjected)
		if err != nil {
			return err
		}

		// reformat interactor._go
		err = utils.Reformat(interactorFilename, interactorBytes)
		if err != nil {
			return err
		}
	}

	return nil

}

func injectCodeToServiceFile(domainName, serviceTemplateCode string) ([]byte, error) {
	return utils.InjectCodeAtTheEndOfFile(getServiceFile(domainName), serviceTemplateCode)
}

// isServiceExist ...
func isServiceExist(servicePath, serviceName string) (bool, error) {
	fset := token.NewFileSet()
	exist := utils.IsExist(fset, servicePath, func(file *ast.File, ts *ast.TypeSpec) bool {
		_, isInterface := ts.Type.(*ast.InterfaceType)
		_, isStruct := ts.Type.(*ast.StructType)
		// TODO we need to handle service as function
		return (isStruct || isInterface) && ts.Name.String() == getServiceName(serviceName)
	})
	return exist, nil
}

func getServiceName(serviceName string) string {
	return fmt.Sprintf("%sService", serviceName)
}

func getServiceInterfaceTemplate() ([]byte, error) {
	return utils.AppTemplates.ReadFile("templates/service/domain_${domainname}/model/service/~service_interface._go")
}

// getServiceInjectTemplate ...
func getServiceInjectTemplate() ([]byte, error) {
	return utils.AppTemplates.ReadFile("templates/service/domain_${domainname}/model/service/~service_inject._go")
}

// injectToOutport ...
func injectToOutport(domainName string, obj ObjTemplate) error {

	fset := token.NewFileSet()

	var iFace *ast.InterfaceType
	var file *ast.File

	rootFolderName := fmt.Sprintf("domain_%s/usecase/%s", domainName, *obj.UsecaseName)
	isAlreadyInjectedBefore := utils.IsExist(fset, rootFolderName, func(f *ast.File, ts *ast.TypeSpec) bool {
		interfaceType, ok := ts.Type.(*ast.InterfaceType)

		if ok && ts.Name.String() == "Outport" {

			file = f
			iFace = interfaceType

			// trace every method to find something line like `service.PublishMessageService`
			for _, meth := range iFace.Methods.List {

				selType, ok := meth.Type.(*ast.SelectorExpr)

				// if interface already injected (PublishMessageService) then abort the mission
				if ok && selType.Sel.String() == getServiceName(obj.ServiceName) {
					return true
				}
			}
		}
		return false
	})

	// if it is already injected, then nothing to do
	if isAlreadyInjectedBefore {
		return nil
	}

	if iFace == nil {
		return fmt.Errorf("outport struct not found")
	}

	// we want to inject it now
	// add new service to outport interface
	iFace.Methods.List = append(iFace.Methods.List, &ast.Field{
		Type: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: utils.GetPackageName(getServiceFolder(domainName)),
			},
			Sel: &ast.Ident{
				Name: getServiceName(obj.ServiceName),
			},
		},
	})

	// TODO who is responsible to write a file? entity or gateway?
	// i prefer to use gateway instead of entity

	fileReadPath := fset.Position(file.Pos()).Filename

	// rewrite the outport
	f, err := os.Create(fileReadPath)
	if err != nil {
		return err
	}

	if err := printer.Fprint(f, fset, file); err != nil {
		return err
	}
	err = f.Close()
	if err != nil {
		return err
	}

	err = utils.Reformat(fileReadPath, nil)
	if err != nil {
		return err
	}

	return nil

}

func getServiceFolder(domainName string) string {
	return fmt.Sprintf("domain_%s/model/service", domainName)
}

func getServiceFile(domainName string) string {
	return fmt.Sprintf("%s/service.go", getServiceFolder(domainName))
}

func getInteractorFilename(domainName, usecaseName string) string {
	return fmt.Sprintf("domain_%s/usecase/%s/interactor.go", domainName, usecaseName)
}

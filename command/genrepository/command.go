package genrepository

import (
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"os"

	"github.com/mirzaakhena/gogen/command/genentity"
	"github.com/mirzaakhena/gogen/utils"
)

// ObjTemplate ...
type ObjTemplate struct {
	PackagePath    string
	RepositoryName string
	EntityName     string
	UsecaseName    *string
}

func Run(inputs ...string) error {

	if len(inputs) < 2 {
		err := fmt.Errorf("\n" +
			"   # Create a repository and inject the template code into interactor file with '//!' flag\n" +
			"   gogen repository SaveOrder Order CreateOrder\n" +
			"     'SaveOrder'   is a repository func name\n" +
			"     'Order'       is an entity name\n" +
			"     'CreateOrder' is an usecase name\n" +
			"\n" +
			"   # Create a repository without inject the template code into usecase\n" +
			"   gogen repository SaveOrder Order\n" +
			"     'SaveOrder' is a repository func name\n" +
			"     'Order'     is an entity name\n" +
			"\n")

		return err
	}

	packagePath := utils.GetPackagePath()
	domainName := utils.GetDefaultDomain()
	repositoryName := inputs[0]
	entityName := inputs[1]

	obj := ObjTemplate{
		PackagePath:    packagePath,
		RepositoryName: repositoryName,
		EntityName:     entityName,
		UsecaseName:    nil,
	}

	if len(inputs) >= 3 {
		obj.UsecaseName = &inputs[2]
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
	}

	err := genentity.Run(entityName)
	if err != nil {
		return err
	}

	err = utils.CreateEverythingExactly("templates/", "repository", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	// repository.go file is already exist, but is the interface is exist ?
	exist, err := isRepoExist(getRepositoryFolder(domainName), repositoryName)
	if err != nil {
		return err
	}

	if !exist {

		// check the prefix and give specific template for it
		templateCode, err := getRepoInterfaceTemplate(obj.RepositoryName)
		if err != nil {
			return err
		}
		templateHasBeenInjected, err := utils.PrintTemplate(string(templateCode), obj)
		if err != nil {
			return err
		}

		dataInBytes, err := injectCodeToRepositoryFile(domainName, templateHasBeenInjected)
		if err != nil {
			return err
		}

		// reformat interactor.go
		err = utils.Reformat(getRepositoryFile(domainName), dataInBytes)
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

		outportFile := getOutportFile(domainName, *obj.UsecaseName)

		// reformat outport._go
		err = utils.Reformat(outportFile, nil)
		if err != nil {
			return err
		}
	}

	// inject to interactor
	{
		// check the prefix and give specific template for it
		interactorCode, err := getRepoInjectTemplate(obj.RepositoryName)
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

func getOutportFile(domainName, usecaseName string) string {
	return fmt.Sprintf("domain_%s/usecase/%s/outport.go", domainName, usecaseName)
}

func injectCodeToRepositoryFile(domainName, repoTemplateCode string) ([]byte, error) {
	filename := getRepositoryFile(domainName)
	return utils.InjectCodeAtTheEndOfFile(filename, repoTemplateCode)
}

func isRepoExist(repoPath, repoName string) (bool, error) {
	fset := token.NewFileSet()
	exist := utils.IsExist(fset, repoPath, func(file *ast.File, ts *ast.TypeSpec) bool {
		_, ok := ts.Type.(*ast.InterfaceType)
		return ok && ts.Name.String() == getRepositoryName(repoName)
	})
	return exist, nil
}

func getRepositoryName(repoName string) string {
	return fmt.Sprintf("%sRepo", repoName)
}

func getRepoInterfaceTemplate(repoName string) ([]byte, error) {

	if utils.HasOneOfThisPrefix(repoName, "save", "create", "add", "update") {
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~repository_interface_save._go")
	}

	if utils.HasOneOfThisPrefix(repoName, "findone", "findfirst", "findlast", "getone") {
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~repository_interface_findone._go")
	}

	if utils.HasOneOfThisPrefix(repoName, "find", "get") {
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~repository_interface_find._go")
	}

	if utils.HasOneOfThisPrefix(repoName, "remove", "delete") {
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~repository_interface_remove._go")
	}

	return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~repository_interface._go")
}

// getRepoInjectTemplate ...
func getRepoInjectTemplate(repoName string) ([]byte, error) {

	if utils.HasOneOfThisPrefix(repoName, "findone", "findfirst", "findlast", "getone") { //
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~interactor_inject_findone._go")
	}

	if utils.HasOneOfThisPrefix(repoName, "find", "get") {
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~interactor_inject_find._go")
	}

	if utils.HasOneOfThisPrefix(repoName, "remove", "delete") {
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~interactor_inject_remove._go")
	}

	if utils.HasOneOfThisPrefix(repoName, "save", "create", "add", "update") {
		return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~interactor_inject_save._go")

	}

	return utils.AppTemplates.ReadFile("templates/repository/domain_${domainname}/model/repository/~interactor_inject._go")

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

			// trace every method to find something line like `repository.SaveOrderRepo`
			for _, meth := range iFace.Methods.List {

				selType, ok := meth.Type.(*ast.SelectorExpr)

				// if interface already injected (SaveOrderRepo) then abort the mission
				if ok && selType.Sel.String() == getRepositoryName(obj.RepositoryName) {
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
	// add new repository to outport interface
	iFace.Methods.List = append(iFace.Methods.List, &ast.Field{
		Type: &ast.SelectorExpr{
			X: &ast.Ident{
				Name: utils.GetPackageName(getRepositoryFolder(domainName)),
			},
			Sel: &ast.Ident{
				Name: getRepositoryName(obj.RepositoryName),
			},
		},
	})

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

func getRepositoryFolder(domainName string) string {
	return fmt.Sprintf("domain_%s/model/repository", domainName)
}

func getRepositoryFile(domainName string) string {
	return fmt.Sprintf("%s/repository.go", getRepositoryFolder(domainName))
}

func getInteractorFilename(domainName, usecaseName string) string {
	return fmt.Sprintf("domain_%s/usecase/%s/interactor.go", domainName, usecaseName)
}

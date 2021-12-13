package entity

import (
  "bufio"
  "bytes"
  "fmt"
  "github.com/mirzaakhena/gogen/domain/domerror"
  "github.com/mirzaakhena/gogen/domain/vo"
  "go/ast"
  "go/printer"
  "go/token"
  "golang.org/x/tools/imports"
  "io/ioutil"
  "os"
  "strings"
)

// ObjRepository ...
type ObjRepository struct {
	RepositoryName vo.Naming
	ObjEntity      ObjEntity
	ObjUsecase     ObjUsecase
}

// ObjDataRepository ...
type ObjDataRepository struct {
	PackagePath    string
	RepositoryName string
	EntityName     string
	UsecaseName    string
}

// NewObjRepository ...
func NewObjRepository(repositoryName, entityName, usecaseName string) (*ObjRepository, error) {

	if repositoryName == "" {
		return nil, domerror.RepositoryNameMustNotEmpty
	}

	var obj ObjRepository
	obj.RepositoryName = vo.Naming(repositoryName)

	{
		et, err := NewObjEntity(entityName)
		if err != nil {
			return nil, err
		}
		obj.ObjEntity = *et
	}

	if usecaseName != "" {
		uc, err := NewObjUsecase(usecaseName)
		if err != nil {
			return nil, err
		}
		obj.ObjUsecase = *uc
	}

	return &obj, nil

}

// GetData ...
func (o ObjRepository) GetData(PackagePath string) *ObjDataRepository {
	return &ObjDataRepository{
		PackagePath:    PackagePath,
		RepositoryName: o.RepositoryName.String(),
		EntityName:     o.ObjEntity.EntityName.String(),
		UsecaseName:    o.ObjUsecase.UsecaseName.String(),
	}
}

// GetRepositoryRootFolderName ...
func (o ObjRepository) GetRepositoryRootFolderName() string {
	return fmt.Sprintf("domain/repository")
}

// GetRepositoryFileName ...
func (o ObjRepository) GetRepositoryFileName() string {
	return fmt.Sprintf("%s/repository.go", o.GetRepositoryRootFolderName())
}

// IsRepoExist ...
func (o ObjRepository) IsRepoExist() (bool, error) {
	fset := token.NewFileSet()
	exist := IsExist(fset, o.GetRepositoryRootFolderName(), func(file *ast.File, ts *ast.TypeSpec) bool {
		_, ok := ts.Type.(*ast.InterfaceType)
		return ok && ts.Name.String() == o.getRepositoryName()
	})
	return exist, nil
}

// InjectCode ...
func (o ObjRepository) InjectCode(repoTemplateCode string) ([]byte, error) {
	return InjectCodeAtTheEndOfFile(o.GetRepositoryFileName(), repoTemplateCode)
}

// InjectToOutport ...
func (o ObjRepository) InjectToOutport() error {

	fset := token.NewFileSet()

	var iFace *ast.InterfaceType
	var file *ast.File

	isAlreadyInjectedBefore := IsExist(fset, o.ObjUsecase.GetUsecaseRootFolderName(), func(f *ast.File, ts *ast.TypeSpec) bool {
		interfaceType, ok := ts.Type.(*ast.InterfaceType)

		if ok && ts.Name.String() == OutportInterfaceName {

			file = f
			iFace = interfaceType

			// trace every method to find something line like `repository.SaveOrderRepo`
			for _, meth := range iFace.Methods.List {

				selType, ok := meth.Type.(*ast.SelectorExpr)

				// if interface already injected (SaveOrderRepo) then abort the mission
				if ok && selType.Sel.String() == o.getRepositoryName() {
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
				Name: GetPackageName(o.GetRepositoryRootFolderName()),
			},
			Sel: &ast.Ident{
				Name: o.getRepositoryName(),
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

	// reformat and import
	newBytes, err := imports.Process(fileReadPath, nil, nil)
	if err != nil {
		return err
	}

	if err := ioutil.WriteFile(fileReadPath, newBytes, 0644); err != nil {
		return err
	}

	return nil

}

// InjectToInteractor ...
func (o ObjRepository) InjectToInteractor(injectedCode string) ([]byte, error) {

	existingFile := o.ObjUsecase.GetInteractorFileName()

	// open interactor file
	file, err := os.Open(existingFile)
	if err != nil {
		return nil, err
	}

	needToInject := false

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		// check the injected code in interactor
		if strings.TrimSpace(row) == injectedCodeLocation {

			needToInject = true

			//// we need to provide an error
			//InitiateError()

			// inject code
			buffer.WriteString(injectedCode)
			buffer.WriteString("\n")

			continue
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	// if no injected marker found, then abort the next step
	if !needToInject {
		return nil, nil
	}

	if err := file.Close(); err != nil {
		return nil, err
	}

	// rewrite the file
	if err := ioutil.WriteFile(existingFile, buffer.Bytes(), 0644); err != nil {
		return nil, err
	}

	return buffer.Bytes(), nil
}

func (o ObjRepository) getRepositoryName() string {
	return fmt.Sprintf("%sRepo", o.RepositoryName)
}

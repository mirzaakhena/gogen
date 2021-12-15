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

// ObjService ...
type ObjService struct {
	ServiceName vo.Naming
	ObjUsecase  ObjUsecase
}

// ObjDataService ...
type ObjDataService struct {
	PackagePath string
	ServiceName string
	UsecaseName string
}

// NewObjService ...
func NewObjService(serviceName, usecaseName string) (*ObjService, error) {

	if serviceName == "" {
		return nil, domerror.ServiceNameMustNotEmpty
	}

	var obj ObjService
	obj.ServiceName = vo.Naming(serviceName)

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
func (o ObjService) GetData(PackagePath string) *ObjDataService {
	return &ObjDataService{
		PackagePath: PackagePath,
		ServiceName: o.ServiceName.String(),
		UsecaseName: o.ObjUsecase.UsecaseName.String(),
	}
}

// GetServiceRootFolderName ...
func (o ObjService) GetServiceRootFolderName() string {
	return fmt.Sprintf("domain/service")
}

// GetServiceFileName ...
func (o ObjService) GetServiceFileName() string {
	return fmt.Sprintf("%s/service.go", o.GetServiceRootFolderName())
}

// IsServiceExist ...
func (o ObjService) IsServiceExist() (bool, error) {
	fset := token.NewFileSet()
	exist := IsExist(fset, o.GetServiceRootFolderName(), func(file *ast.File, ts *ast.TypeSpec) bool {
		_, isInterface := ts.Type.(*ast.InterfaceType)
		_, isStruct := ts.Type.(*ast.StructType)
		// TODO we need to handle service as function
		return (isStruct || isInterface) && ts.Name.String() == o.getServiceName()
	})
	return exist, nil
}

// InjectCode ...
func (o ObjService) InjectCode(repoTemplateCode string) ([]byte, error) {
	return InjectCodeAtTheEndOfFile(o.GetServiceFileName(), repoTemplateCode)
}

// InjectToOutport ...
func (o ObjService) InjectToOutport() error {

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
				if ok && selType.Sel.String() == o.getServiceName() {
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
				Name: GetPackageName(o.GetServiceRootFolderName()),
			},
			Sel: &ast.Ident{
				Name: o.getServiceName(),
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
func (o ObjService) InjectToInteractor(injectedCode string) ([]byte, error) {

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

func (o ObjService) getServiceName() string {
	return fmt.Sprintf("%sService", o.ServiceName)
}

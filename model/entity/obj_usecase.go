package entity

import (
	"fmt"
	"github.com/mirzaakhena/gogen/model/domerror"
	"github.com/mirzaakhena/gogen/model/vo"
)

const (
	// OutportInterfaceName ...
	OutportInterfaceName = "Outport"
)

// ObjUsecase ...
type ObjUsecase struct {
	UsecaseName vo.Naming
}

// ObjDataUsecase ...
type ObjDataUsecase struct {
	PackagePath string
	UsecaseName string
}

// NewObjUsecase ...
func NewObjUsecase(usecaseName string) (*ObjUsecase, error) {

	if usecaseName == "" {
		return nil, domerror.UsecaseNameMustNotEmpty
	}

	var obj ObjUsecase
	obj.UsecaseName = vo.Naming(usecaseName)

	return &obj, nil
}

// GetData ...
func (o ObjUsecase) GetData(PackagePath string) *ObjDataUsecase {
	return &ObjDataUsecase{
		PackagePath: PackagePath,
		UsecaseName: o.UsecaseName.String(),
	}
}

// TODO create new usecase will create new interactor

// GetUsecaseRootFolderName ...
func (o ObjUsecase) GetUsecaseRootFolderName() string {
	return fmt.Sprintf("usecase/%s", o.UsecaseName.LowerCase())
}

// GetInportFileName ...
func (o ObjUsecase) GetInportFileName() string {
	return fmt.Sprintf("%s/inport.go", o.GetUsecaseRootFolderName())
}

// GetOutportFileName ...
func (o ObjUsecase) GetOutportFileName() string {
	return fmt.Sprintf("%s/outport.go", o.GetUsecaseRootFolderName())
}

// GetInteractorFileName ...
func (o ObjUsecase) GetInteractorFileName() string {
	return fmt.Sprintf("%s/interactor.go", o.GetUsecaseRootFolderName())
}

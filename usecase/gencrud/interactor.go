package gencrud

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
)

//go:generate mockery --name Outport -output mocks/

type genCrudInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(outputPort Outport) Inport {
	return &genCrudInteractor{
		outport: outputPort,
	}
}

// Execute the usecase
func (r *genCrudInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// create object service
	obj, err := entity.NewObjCrud(req.EntityName)
	if err != nil {
		return nil, err
	}

	packagePath := r.outport.GetPackagePath(ctx)

	filerenamer := map[string]string{
		"entityname": obj.Entity.EntityName.LowerCase(),
	}

	err = service.CreateEverythingExactly("crud/", "application", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("crud/", "controller", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("crud/", "domain", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("crud/", "gateway", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("crud/", "infrastructure", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("crud/", "usecase", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	{
		templateFile := r.outport.GetMainFileForCrudTemplate(ctx)
		_, err = r.outport.WriteFileIfNotExist(ctx, templateFile, "main.go", obj.GetData(packagePath))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

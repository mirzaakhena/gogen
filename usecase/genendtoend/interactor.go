package genendtoend

import (
	"context"
	"github.com/mirzaakhena/gogen/model/entity"
	"github.com/mirzaakhena/gogen/model/service"
)

//go:generate mockery --name Outport -output mocks/

type genEndToEndInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(outputPort Outport) Inport {
	return &genEndToEndInteractor{
		outport: outputPort,
	}
}

// Execute the usecase
func (r *genEndToEndInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// create object service
	obj, err := entity.NewObjEndToEnd(req.UsecaseName, req.EntityName)
	if err != nil {
		return nil, err
	}

	packagePath := r.outport.GetPackagePath(ctx)

	filerenamer := map[string]string{
		"usecasename": obj.Usecase.UsecaseName.LowerCase(),
		"entityname":  obj.Entity.EntityName.LowerCase(),
	}

	err = service.CreateEverythingExactly("endtoend/", "application", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("endtoend/", "controller", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("endtoend/", "model", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("endtoend/", "gateway", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("endtoend/", "infrastructure", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("endtoend/", "usecase", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	{
		templateFile := r.outport.GetMainFileForE2ETemplate(ctx)
		_, err = r.outport.WriteFileIfNotExist(ctx, templateFile, "main.go", obj.GetData(packagePath))
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

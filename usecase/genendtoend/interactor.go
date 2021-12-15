package genendtoend

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
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

	err = service.CreateEverythingExactly("endtoend/", "endtoend", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	return res, nil
}

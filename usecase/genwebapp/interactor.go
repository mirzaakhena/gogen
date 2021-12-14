package genwebapp

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
)

//go:generate mockery --name Outport -output mocks/

type genWebappInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(outputPort Outport) Inport {
	return &genWebappInteractor{
		outport: outputPort,
	}
}

// Execute the usecase
func (r *genWebappInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

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

	err = service.CreateEverythingExactly("webapp/", "webapp", filerenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	return res, nil
}

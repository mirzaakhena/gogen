package genvaluestring

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/domain/entity"
)

//go:generate mockery --name Outport -output mocks/

type genValueStringInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenValueString
func NewUsecase(outputPort Outport) Inport {
	return &genValueStringInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenValueString
func (r *genValueStringInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	obj, err := entity.NewObjValueString(req.ValueStringName)
	if err != nil {
		return nil, err
	}

	// code your usecase definition here ...
	_, err = r.outport.CreateFolderIfNotExist(ctx,"domain/vo")
	if err != nil {
		return nil, err
	}

	{
		outputFile := fmt.Sprintf("domain/vo/%s.go", obj.ValueStringName.SnakeCase())
		tem := r.outport.GetValueStringTemplate(ctx)
		_, err = r.outport.WriteFileIfNotExist(ctx, tem, outputFile, obj.GetData())
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

package genvalueobject

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/domain/entity"
)

//go:generate mockery --name Outport -output mocks/

type genValueObjectInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenValueObject
func NewUsecase(outputPort Outport) Inport {
	return &genValueObjectInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenValueObject
func (r *genValueObjectInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	obj, err := entity.NewObjValueObject(req.ValueObjectName, req.FieldNames)
	if err != nil {
		return nil, err
	}

	// code your usecase definition here ...
	_, err = r.outport.CreateFolderIfNotExist(ctx,"domain/vo")
	if err != nil {
		return nil, err
	}

	{
		outputFile := fmt.Sprintf("domain/vo/%s.go", obj.ValueObjectName.SnakeCase())
		tem := r.outport.GetValueObjectTemplate(ctx)
		_, err = r.outport.WriteFileIfNotExist(ctx, tem, outputFile, obj.GetData())
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

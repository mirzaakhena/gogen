package generror

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
)

//go:generate mockery --name Outport -output mocks/

type genErrorInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenError
func NewUsecase(outputPort Outport) Inport {
	return &genErrorInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenError
func (r *genErrorInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	//packagePath := r.outport.GetPackagePath(ctx)

	err := service.CreateEverythingExactly("default/", "domain/domerror", map[string]string{}, struct{}{})
	if err != nil {
		return nil, err
	}

	objError, err := entity.NewObjError(req.ErrorName)
	if err != nil {
		return nil, err
	}

	errorLine := r.outport.GetErrorLineTemplate(ctx)

	// we will only inject the non existing method
	data := objError.GetData()

	templateHasBeenInjected, err := r.outport.PrintTemplate(ctx, errorLine, data)
	if err != nil {
		return nil, err
	}

	bytes, err := objError.InjectCode(templateHasBeenInjected)
	if err != nil {
		return nil, err
	}

	// reformat outport._go
	err = r.outport.Reformat(ctx, objError.GetErrorEnumFileName(), bytes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

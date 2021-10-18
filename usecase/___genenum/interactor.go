package ___genenum

import "context"

//go:generate mockery --name Outport -output mocks/

type genEnumInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenEnum
func NewUsecase(outputPort Outport) Inport {
	return &genEnumInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenEnum
func (r *genEnumInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// code your usecase definition here ...

	return res, nil
}

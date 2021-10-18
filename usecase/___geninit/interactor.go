package ___geninit

import "context"

//go:generate mockery --name Outport -output mocks/

type genInitInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenInit
func NewUsecase(outputPort Outport) Inport {
	return &genInitInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenInit
func (r *genInitInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// code your usecase definition here ...

	return res, nil
}

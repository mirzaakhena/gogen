package gethello

import (
	"context"
)

//go:generate mockery --name Outport -output mocks/

type getHelloInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase
func NewUsecase(outputPort Outport) Inport {
	return &getHelloInteractor{
		outport: outputPort,
	}
}

// Execute the usecase
func (r *getHelloInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	res.Message = "Hello World"

	return res, nil
}

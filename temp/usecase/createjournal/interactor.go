package createjournal

import (
	"context"
)

//go:generate mockery --name Outport -output mocks/

type createJournalInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase CreateJournal
func NewUsecase(outputPort Outport) Inport {
	return &createJournalInteractor{
		outport: outputPort,
	}
}

// Execute the usecase CreateJournal
func (r *createJournalInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	return res, nil
}

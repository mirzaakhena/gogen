package readjournal

import "context"

//go:generate mockery --dir port/ --name ReadJournalOutport -output mocks/

type readJournalInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase ReadJournal
func NewUsecase(outputPort Outport) Inport {
	return &readJournalInteractor{
		outport: outputPort,
	}
}

// Execute the usecase ReadJournal
func (r *readJournalInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// code your usecase definition here ...

	return res, nil
}

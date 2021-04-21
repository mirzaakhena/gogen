package readjournal

import (
	"context"
)

// Inport of ReadJournal
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase ReadJournal
type InportRequest struct {
}

// InportResponse is response payload after running the usecase ReadJournal
type InportResponse struct {
}

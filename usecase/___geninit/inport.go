package ___geninit

import (
	"context"
)

// Inport of GenInit
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenInit
type InportRequest struct {
}

// InportResponse is response payload after running the usecase GenInit
type InportResponse struct {
}

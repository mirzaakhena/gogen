package ___genenum

import (
	"context"
)

// Inport of GenEnum
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenEnum
type InportRequest struct {
}

// InportResponse is response payload after running the usecase GenEnum
type InportResponse struct {
}

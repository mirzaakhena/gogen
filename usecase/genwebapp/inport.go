package genwebapp

import (
	"context"
)

// Inport of Usecase
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase
type InportRequest struct {
	EntityName string
}

// InportResponse is response payload after running the usecase
type InportResponse struct {
}

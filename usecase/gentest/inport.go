package gentest

import (
	"context"
)

// Inport of GenTest
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenTest
type InportRequest struct {
	UsecaseName string
	TestName    string
}

// InportResponse is response payload after running the usecase GenTest
type InportResponse struct {
	Message string
}

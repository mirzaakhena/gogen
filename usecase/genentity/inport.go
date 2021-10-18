package genentity

import (
	"context"
)

// Inport of GenEntity
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenEntity
type InportRequest struct {
	EntityName string
}

// InportResponse is response payload after running the usecase GenEntity
type InportResponse struct {
	Message string
}

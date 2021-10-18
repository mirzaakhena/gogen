package generror

import (
	"context"
)

// Inport of GenError
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenError
type InportRequest struct {
	ErrorName string
}

// InportResponse is response payload after running the usecase GenError
type InportResponse struct {
}

package genusecase

import (
	"context"
)

// Inport of GenUsecase
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenUsecase
type InportRequest struct {
	UsecaseName string
}

// InportResponse is response payload after running the usecase GenUsecase
type InportResponse struct {
	Message string
}

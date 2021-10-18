package genrepository

import (
	"context"
)

// Inport of GenRepository
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenRepository
type InportRequest struct {
	RepositoryName string
	EntityName     string
	UsecaseName    string
}

// InportResponse is response payload after running the usecase GenRepository
type InportResponse struct {
	Message string
}

package gengateway

import (
	"context"
)

// Inport of GenGateway
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenGateway
type InportRequest struct {
	UsecaseName string //
	GatewayName string //
}

// InportResponse is response payload after running the usecase GenGateway
type InportResponse struct {
}

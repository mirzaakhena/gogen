package genregistry

import (
  "context"
)

// Inport of GenRegistry
type Inport interface {
  Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenRegistry
type InportRequest struct {
  RegistryName   string
  ControllerName string
  GatewayName    string
}

// InportResponse is response payload after running the usecase GenRegistry
type InportResponse struct {
  Message string
}

package gencontroller

import (
  "context"
)

// Inport of GenController
type Inport interface {
  Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenController
type InportRequest struct {
  ControllerName string
  UsecaseName    string
  DriverName     string
}

// InportResponse is response payload after running the usecase GenController
type InportResponse struct {
}

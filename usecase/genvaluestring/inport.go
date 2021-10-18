package genvaluestring

import (
  "context"
)

// Inport of GenValueString
type Inport interface {
  Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenValueString
type InportRequest struct {
  ValueStringName string
}

// InportResponse is response payload after running the usecase GenValueString
type InportResponse struct {
}

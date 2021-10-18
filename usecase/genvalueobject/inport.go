package genvalueobject

import (
  "context"
)

// Inport of GenValueObject
type Inport interface {
  Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase GenValueObject
type InportRequest struct {
  ValueObjectName string
  FieldNames      []string
}

// InportResponse is response payload after running the usecase GenValueObject
type InportResponse struct {
}

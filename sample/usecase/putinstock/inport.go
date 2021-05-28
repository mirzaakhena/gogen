package putinstock

import (
	"context"
	"time"
)

// Inport of PutInStock
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase PutInStock
type InportRequest struct {
	Date          time.Time
	Quantity      int
	TotalPrice    float64
	InventoryCode string
	Method        string
}

// InportResponse is response payload after running the usecase PutInStock
type InportResponse struct {
}

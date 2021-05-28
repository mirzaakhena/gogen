package pulloutstock

import (
	"context"
	"time"
)

// Inport of PullOutStock
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase PullOutStock
type InportRequest struct {
	Date          time.Time
	Quantity      int
	InventoryCode string
}

// InportResponse is response payload after running the usecase PullOutStock
type InportResponse struct {
}

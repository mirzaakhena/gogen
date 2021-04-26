package updatestock

import (
	"context"
	"time"
)

// Inport of UpdateStock
type Inport interface {
	Execute(ctx context.Context, req InportRequest) (*InportResponse, error)
}

// InportRequest is request payload to run the usecase UpdateStock
type InportRequest struct {
	Date          time.Time `` //
	ReferenceID   string    `` //
	Description   string    `` //
	InventoryCode string    `` //
	Price         float64   `` //
	Quantity      int       `` //
}

// InportResponse is response payload after running the usecase UpdateStock
type InportResponse struct {
}

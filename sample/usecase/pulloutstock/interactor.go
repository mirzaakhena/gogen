package pulloutstock

import (
	"accounting/domain/entity"
	"context"
)

//go:generate mockery --name Outport -output mocks/

type pullOutStockInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase PullOutStock
func NewUsecase(outputPort Outport) Inport {
	return &pullOutStockInteractor{
		outport: outputPort,
	}
}

// Execute the usecase PullOutStock
func (r *pullOutStockInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	lastStockPriceObj, err := r.outport.FindLastStockPrice(ctx, req.InventoryCode)
	if err != nil {
		return nil, err
	}

	stockPriceObj, err := entity.NewStockPrice(entity.StockPriceRequest{
		GetID: func() string {
			return r.outport.GenerateUUID(ctx)
		},
		Date:          req.Date,
		InventoryCode: req.InventoryCode,
	})

	err = stockPriceObj.PullOut(req.Quantity, lastStockPriceObj)
	if err != nil {
		return nil, err
	}

	if lastStockPriceObj != nil && lastStockPriceObj.UpdatedByID != "" {
		err = r.outport.SaveStockPrice(ctx, lastStockPriceObj)
		if err != nil {
			return nil, err
		}
	}

	err = r.outport.SaveStockPrice(ctx, stockPriceObj)
	if err != nil {
		return nil, err
	}

	return res, nil
}

package putinstock

import (
	"accounting/domain/entity"
	"context"
)

//go:generate mockery --name Outport -output mocks/

type putInStockInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase PutInStock
func NewUsecase(outputPort Outport) Inport {
	return &putInStockInteractor{
		outport: outputPort,
	}
}

// Execute the usecase PutInStock
func (r *putInStockInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

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

	err = stockPriceObj.PutIn(req.Quantity, req.TotalPrice, lastStockPriceObj, req.Method)
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

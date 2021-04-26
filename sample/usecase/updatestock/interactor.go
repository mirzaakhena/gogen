package updatestock

import (
	"accounting/domain/entity"
	"context"
)

//go:generate mockery --name Outport -output mocks/

type updateStockInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase UpdateStock
func NewUsecase(outputPort Outport) Inport {
	return &updateStockInteractor{
		outport: outputPort,
	}
}

// Execute the usecase UpdateStock
func (r *updateStockInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	inventoryStockObjs, err := r.outport.FindLastQuantityAndPrice(ctx, req.InventoryCode)
	if err != nil {
		return nil, err
	}

	lastBalancePrice, err := inventoryStockObjs.BalancePrice()
	if err != nil {
		return nil, err
	}

	inventoryStockObj, err := entity.NewInventoryStock(entity.InventoryStockRequest{
		Date:                req.Date,
		ReferenceID:         req.ReferenceID,
		Description:         req.Description,
		InventoryCode:       req.InventoryCode,
		Price:               req.Price,
		Quantity:            req.Quantity,
		LastBalanceQuantity: inventoryStockObjs.BalanceQuantity,
		LastBalancePrice:    lastBalancePrice,
	})
	if err != nil {
		return nil, err
	}

	err = r.outport.SaveInventoryStock(ctx, inventoryStockObj)
	if err != nil {
		return nil, err
	}

	return res, nil
}

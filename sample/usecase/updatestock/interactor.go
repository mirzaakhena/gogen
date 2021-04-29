package updatestock

import (
	"accounting/application/apperror"
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

	inventoryObj, err := r.outport.FindOneInventory(ctx, req.InventoryCode)
	if err != nil {
		return nil, err
	}
	if inventoryObj == nil {
		return nil, apperror.ObjectNotFound.Var(inventoryObj)
	}

	stockPriceObjs, err := r.outport.FindLastStockPrice(ctx, req.InventoryCode)
	if err != nil {
		return nil, err
	}

	inventoryDataObj, err := entity.NewInventoryBalance(entity.InventoryBalanceRequest{
		Date:              req.Date,
		ReferenceID:       req.ReferenceID,
		Description:       req.Description,
		InventoryCode:     req.InventoryCode,
		CalculationMethod: inventoryObj.CalculationMethod,
		TotalPrice:        req.TotalPrice,
		Quantity:          req.Quantity,
		LastStockPrices:   stockPriceObjs,
	})
	if err != nil {
		return nil, err
	}

	err = r.outport.SaveInventoryStock(ctx, inventoryDataObj)
	if err != nil {
		return nil, err
	}

	return res, nil
}

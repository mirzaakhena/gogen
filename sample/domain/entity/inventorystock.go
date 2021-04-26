package entity

import (
	"accounting/application/apperror"
	"time"
)

type InventoryStock struct {
	ID                string    `` //
	Date              time.Time `` //
	ReferenceID       string    `` //
	Description       string    `` //
	InventoryCode     string    `` //
	AmountQuantity    int       `` //
	AmountPrice       float64   `` //
	BalanceQuantity   int       `` //
	BalanceTotalPrice float64   `` //
}

type InventoryStockRequest struct {
	Date                time.Time `` //
	ReferenceID         string    `` //
	Description         string    `` //
	InventoryCode       string    `` //
	Price               float64   `` //
	Quantity            int       `` //
	LastBalanceQuantity int       `` //
	LastBalancePrice    float64   `` //
}

func NewInventoryStock(req InventoryStockRequest) (*InventoryStock, error) {

	if req.LastBalanceQuantity+req.Quantity < 0 {
		return nil, apperror.NotEnoughSufficientQuantity
	}

	if req.Quantity > 0 && req.Price < 0 {
		return nil, apperror.PriceMustNotBelowZero
	}

	var obj InventoryStock
	obj.ReferenceID = req.ReferenceID
	obj.Description = req.Description
	obj.Date = req.Date
	obj.InventoryCode = req.InventoryCode
	obj.AmountQuantity = req.Quantity
	obj.AmountPrice = req.Price
	obj.BalanceQuantity = req.LastBalanceQuantity + req.Quantity
	obj.BalanceTotalPrice = req.LastBalancePrice + req.Price

	return &obj, nil
}

func (i *InventoryStock) AmountTotalPrice() float64 {
	return i.AmountPrice * float64(i.AmountQuantity)
}

func (i *InventoryStock) BalancePrice() (float64, error) {
	if i.BalanceQuantity <= 0 {
		return 0, apperror.StockQuantityIsZero
	}
	return i.BalanceTotalPrice / float64(i.BalanceQuantity), nil
}

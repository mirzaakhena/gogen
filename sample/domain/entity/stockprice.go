package entity

import (
	"accounting/application/apperror"
	"math"
	"time"
)

type StockPrice struct {
	ID            string `` //
	Date          time.Time
	Quantity      int
	Price         float64
	InventoryCode string
	UpdatedByID   string
}

type StockPriceRequest struct {
	GetID         func() string
	Date          time.Time
	InventoryCode string
}

func NewStockPrice(req StockPriceRequest) (*StockPrice, error) {

	if req.InventoryCode == "" {
		return nil, apperror.InventoryCodeMustNotEmpty
	}

	generatedID := req.GetID()

	if generatedID == "" {
		return nil, apperror.InventoryCodeMustNotEmpty
	}

	var obj StockPrice
	obj.ID = generatedID
	obj.Date = req.Date
	obj.InventoryCode = req.InventoryCode

	return &obj, nil

}

func (r *StockPrice) PutIn(qty int, totalPrice float64, lastStockPrice *StockPrice, method string) error {

	if qty <= 0 {
		return apperror.QuantityMustNotZero
	}

	if lastStockPrice == nil {

		r.Quantity = qty
		pricePerQty := totalPrice / float64(qty)
		r.Price = pricePerQty

		return nil
	}

	if method == "WAVG" {

		newQuantity := lastStockPrice.Quantity + qty
		newTotalPrice := lastStockPrice.GetTotalPrice() + totalPrice
		r.Quantity = newQuantity
		r.Price = math.Floor(newTotalPrice / float64(newQuantity))

		// invalidate last stock price
		lastStockPrice.UpdatedByID = r.ID

		return nil
	}

	pricePerQty := math.Floor(totalPrice / float64(qty))
	r.Price = pricePerQty

	if lastStockPrice.Price == r.Price {
		r.Quantity = lastStockPrice.Quantity + qty

		// invalidate last stock price
		lastStockPrice.UpdatedByID = r.ID

		return nil
	}

	r.Quantity = qty

	return nil

}

func (r *StockPrice) PullOut(qty int, lastStockPrice *StockPrice) error {

	if lastStockPrice == nil {
		return apperror.StockIsEmpty
	}

	qtyReminder := lastStockPrice.Quantity - qty
	if qtyReminder < 0 {
		return apperror.StockIsNotEnough
	}

	r.Quantity = qtyReminder
	r.Price = lastStockPrice.Price

	// invalidate last stock price
	lastStockPrice.UpdatedByID = r.ID

	return nil
}

func (r *StockPrice) GetTotalPrice() float64 {
	return r.Price * float64(r.Quantity)
}

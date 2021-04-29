package entity

import (
	"accounting/application/apperror"
	"time"
)

type InventoryBalance struct {
	ID            string        `` //
	Date          time.Time     `` //
	ReferenceID   string        `` //
	Description   string        `` //
	InventoryCode string        `` //
	Quantity      int           `` //
	Price         float64       `` //
	StockPrices   []*StockPrice `` //
}

type InventoryBalanceRequest struct {
	Date              time.Time     `` //
	ReferenceID       string        `` //
	Description       string        `` //
	CalculationMethod string        `` //
	InventoryCode     string        `` //
	TotalPrice        float64       `` //
	Quantity          int           `` //
	LastStockPrices   []*StockPrice `` //
}

func NewInventoryBalance(req InventoryBalanceRequest) (*InventoryBalance, error) {

	// for put new stock, the price must not negative
	if req.Quantity > 0 && req.TotalPrice < 0 {
		return nil, apperror.PriceMustNotBelowZero
	}

	newStockPrices := make([]*StockPrice, 0)

	// no stock yet
	if len(req.LastStockPrices) == 0 {

		// pull out will produce : no stock available
		if req.Quantity < 0 {
			return nil, apperror.NotEnoughSufficientQuantity
		}

		// put in
		// will inserting the stock

		pricePerOneQty := req.TotalPrice / float64(req.Quantity)

		newStockPrices = append(newStockPrices, &StockPrice{
			Quantity: req.Quantity,
			Price:    pricePerOneQty,
		})

	} else // we have existing stock

	{

		if req.CalculationMethod == "WAVG" {

			// pull out
			if req.Quantity < 0 {

				// check the available stock
				remainingQuantity := req.LastStockPrices[0].Quantity + req.Quantity

				// if not available then stop
				if remainingQuantity < 0 {
					return nil, apperror.NotEnoughSufficientQuantity
				}

				// otherwise will continue
				usedPrice := req.LastStockPrices[0].Price

				newStockPrices = append(newStockPrices, &StockPrice{
					Quantity: remainingQuantity,
					Price:    usedPrice,
				})

			} else

			// put in
			{

				newQuantity := req.LastStockPrices[0].Quantity + req.Quantity

				newBalanceTotalPrice := req.LastStockPrices[0].GetTotalPrice() + req.TotalPrice

				pricePerOneQty := newBalanceTotalPrice / float64(newQuantity)

				newStockPrices = append(newStockPrices, &StockPrice{
					Quantity: newQuantity,
					Price:    pricePerOneQty,
				})

			}

		} else //

		if req.CalculationMethod == "FIFO" {

			// pull out
			if req.Quantity < 0 {

				// we store the absolute value of Qty
				requestedQuantity := -req.Quantity

				for i := 0; i < len(req.LastStockPrices); i++ {

					sp := req.LastStockPrices[i]

					if requestedQuantity <= sp.Quantity {
						sp.Quantity -= requestedQuantity
						break
					} else {
						requestedQuantity -= sp.Quantity
						sp.Quantity = 0
					}

				}

				// restore the stock
				for i := 0; i < len(req.LastStockPrices); i++ {

					sp := req.LastStockPrices[i]

					if sp.Quantity == 0 {
						continue
					}

					newStockPrices = append(newStockPrices, &StockPrice{
						Quantity: sp.Quantity,
						Price:    sp.Price,
					})

				}

			} else

			// put in
			{
				sp := req.LastStockPrices[0]

				pricePerQty := req.TotalPrice / float64(req.Quantity)
				if sp.Price == pricePerQty {
					sp.Quantity += req.Quantity

					newStockPrices = req.LastStockPrices

				} else //

				{

					newStockPrices = req.LastStockPrices

					newStockPrices = append(newStockPrices, &StockPrice{
						Quantity: req.Quantity,
						Price:    sp.Price,
					})

				}

			}

		} else //

		if req.CalculationMethod == "LIFO" {

			// pull out
			if req.Quantity < 0 {

				// we store the absolute value of Qty
				requestedQuantity := -req.Quantity

				n := len(req.LastStockPrices)
				for i := 0; i < n; i++ {

					sp := req.LastStockPrices[n-i-1]

					if requestedQuantity <= sp.Quantity {
						sp.Quantity -= requestedQuantity
						break
					} else {
						requestedQuantity -= sp.Quantity
						sp.Quantity = 0
					}

				}

				// restore the stock
				for i := 0; i < len(req.LastStockPrices); i++ {

					sp := req.LastStockPrices[i]

					if sp.Quantity == 0 {
						continue
					}

					newStockPrices = append(newStockPrices, &StockPrice{
						Quantity: sp.Quantity,
						Price:    sp.Price,
					})

				}

			} else //

			{

				n := len(req.LastStockPrices)

				sp := req.LastStockPrices[n-1]

				pricePerQty := req.TotalPrice / float64(req.Quantity)
				if sp.Price == pricePerQty {
					sp.Quantity += req.Quantity

					newStockPrices = req.LastStockPrices

				} else //

				{

					newStockPrices = req.LastStockPrices

					newStockPrices = append(newStockPrices, &StockPrice{
						Quantity: req.Quantity,
						Price:    sp.Price,
					})

				}
			}

		}

	}

	var obj InventoryBalance
	obj.ReferenceID = req.ReferenceID
	obj.Description = req.Description
	obj.Date = req.Date
	obj.InventoryCode = req.InventoryCode
	obj.Quantity = req.Quantity
	obj.Price = req.TotalPrice / float64(req.Quantity)
	obj.StockPrices = newStockPrices

	return &obj, nil
}

func (i *InventoryBalance) AmountTotalPrice() float64 {
	return i.Price * float64(i.Quantity)
}

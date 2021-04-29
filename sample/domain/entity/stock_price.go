package entity

type StockPrice struct {
	Quantity int     `` //
	Price    float64 `` //
}

type StockPriceRequest struct {
}

func NewStockPrice(req StockPriceRequest) (*StockPrice, error) {

	var obj StockPrice

	return &obj, nil
}

func (r StockPrice) GetTotalPrice() float64 {
	return r.Price * float64(r.Quantity)
}

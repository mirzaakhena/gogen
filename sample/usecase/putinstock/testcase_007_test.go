package putinstock

import (
	"accounting/domain/entity"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// FIFO
//
//input
//qty : 2
//totalPrice : 4200
//
//lastStock
//qty : 5
//price : 2000
//
//result
//qty : 2
//price : 2100

type mockOutport7 struct {
	t         *testing.T
	nowDate   time.Time
	lastQty   int
	lastPrice float64
}

func (m mockOutport7) SaveStockPrice(ctx context.Context, obj *entity.StockPrice) error {

	assert.Equal(m.t, obj.ID, "35")
	assert.Equal(m.t, obj.InventoryCode, "X12")
	assert.Equal(m.t, obj.UpdatedByID, "")
	assert.Equal(m.t, obj.Quantity, 2)
	assert.Equal(m.t, obj.Price, float64(2100))
	assert.Equal(m.t, obj.Date, m.nowDate)

	return nil
}

func (m mockOutport7) FindLastStockPrice(ctx context.Context, stockPriceID string) (*entity.StockPrice, error) {

	return &entity.StockPrice{
		ID:            "34",
		Date:          m.nowDate.Add(-1 * time.Hour),
		Quantity:      m.lastQty,
		Price:         m.lastPrice,
		InventoryCode: "X12",
		UpdatedByID:   "",
	}, nil

}

func (m mockOutport7) GenerateUUID(ctx context.Context) string {
	return "35"
}

func TestNormal7(t *testing.T) {

	now := time.Now()

	mockOutport7Obj := mockOutport7{
		t:         t,
		nowDate:   now,
		lastQty:   5,
		lastPrice: 2000,
	}

	ctx := context.Background()

	_, err := NewUsecase(&mockOutport7Obj).Execute(ctx, InportRequest{
		Date:          now,
		Quantity:      2,
		TotalPrice:    4200,
		InventoryCode: "X12",
		Method:        "FIFO",
	})

	assert.Nil(t, err)

}

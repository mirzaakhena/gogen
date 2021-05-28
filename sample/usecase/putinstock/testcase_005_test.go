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
//totalPrice : 4000
//
//lastStock
//qty : 5
//price : 2000
//
//result
//qty : 7
//price : 2000

type mockOutport5 struct {
	t         *testing.T
	nowDate   time.Time
	lastQty   int
	lastPrice float64
}

func (m mockOutport5) SaveStockPrice(ctx context.Context, obj *entity.StockPrice) error {

	if obj.ID == "34" {

		assert.Equal(m.t, obj.UpdatedByID, "35")
		assert.Equal(m.t, obj.InventoryCode, "X12")
		assert.Equal(m.t, obj.Quantity, 5)
		assert.Equal(m.t, obj.Price, float64(2000))
		assert.Equal(m.t, obj.Date, m.nowDate.Add(-1*time.Hour))

		return nil
	}

	assert.Equal(m.t, obj.ID, "35")
	assert.Equal(m.t, obj.InventoryCode, "X12")
	assert.Equal(m.t, obj.UpdatedByID, "")
	assert.Equal(m.t, obj.Quantity, 7)
	assert.Equal(m.t, obj.Price, float64(2000))
	assert.Equal(m.t, obj.Date, m.nowDate)

	return nil
}

func (m mockOutport5) FindLastStockPrice(ctx context.Context, stockPriceID string) (*entity.StockPrice, error) {

	return &entity.StockPrice{
		ID:            "34",
		Date:          m.nowDate.Add(-1 * time.Hour),
		Quantity:      m.lastQty,
		Price:         m.lastPrice,
		InventoryCode: "X12",
		UpdatedByID:   "",
	}, nil

}

func (m mockOutport5) GenerateUUID(ctx context.Context) string {
	return "35"
}

func TestNormal5(t *testing.T) {

	now := time.Now()

	mockOutport5Obj := mockOutport5{
		t:         t,
		nowDate:   now,
		lastQty:   5,
		lastPrice: 2000,
	}

	ctx := context.Background()

	_, err := NewUsecase(&mockOutport5Obj).Execute(ctx, InportRequest{
		Date:          now,
		Quantity:      2,
		TotalPrice:    4000,
		InventoryCode: "X12",
		Method:        "FIFO",
	})

	assert.Nil(t, err)

}

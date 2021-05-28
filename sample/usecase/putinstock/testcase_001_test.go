package putinstock

import (
	"accounting/domain/entity"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// WAVG
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

type mockOutport1 struct {
	t         *testing.T
	nowDate   time.Time
	lastQty   int
	lastPrice float64
}

func (m mockOutport1) SaveStockPrice(ctx context.Context, obj *entity.StockPrice) error {

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

func (m mockOutport1) FindLastStockPrice(ctx context.Context, stockPriceID string) (*entity.StockPrice, error) {

	return &entity.StockPrice{
		ID:            "34",
		Date:          m.nowDate.Add(-1 * time.Hour),
		Quantity:      m.lastQty,
		Price:         m.lastPrice,
		InventoryCode: "X12",
		UpdatedByID:   "",
	}, nil

}

func (m mockOutport1) GenerateUUID(ctx context.Context) string {
	return "35"
}

func TestNormal1(t *testing.T) {

	now := time.Now()

	mockOutport1Obj := mockOutport1{
		t:         t,
		nowDate:   now,
		lastQty:   5,
		lastPrice: 2000,
	}

	ctx := context.Background()

	_, err := NewUsecase(&mockOutport1Obj).Execute(ctx, InportRequest{
		Date:          now,
		Quantity:      2,
		TotalPrice:    4000,
		InventoryCode: "X12",
		Method:        "WAVG",
	})

	assert.Nil(t, err)

}

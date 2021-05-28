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
//qty : 0
//price : 2000
//
//result
//qty : 2
//price : 2000

type mockOutport6 struct {
	t         *testing.T
	nowDate   time.Time
	lastQty   int
	lastPrice float64
}

func (m mockOutport6) SaveStockPrice(ctx context.Context, obj *entity.StockPrice) error {

	assert.Equal(m.t, obj.ID, "35")
	assert.Equal(m.t, obj.InventoryCode, "X12")
	assert.Equal(m.t, obj.UpdatedByID, "")
	assert.Equal(m.t, obj.Quantity, 2)
	assert.Equal(m.t, obj.Price, float64(2000))
	assert.Equal(m.t, obj.Date, m.nowDate)

	return nil
}

func (m mockOutport6) FindLastStockPrice(ctx context.Context, stockPriceID string) (*entity.StockPrice, error) {

	return nil, nil

}

func (m mockOutport6) GenerateUUID(ctx context.Context) string {
	return "35"
}

func TestNormal6(t *testing.T) {

	now := time.Now()

	mockOutport6Obj := mockOutport6{
		t:         t,
		nowDate:   now,
		lastQty:   5,
		lastPrice: 2000,
	}

	ctx := context.Background()

	_, err := NewUsecase(&mockOutport6Obj).Execute(ctx, InportRequest{
		Date:          now,
		Quantity:      2,
		TotalPrice:    4000,
		InventoryCode: "X12",
		Method:        "FIFO",
	})

	assert.Nil(t, err)

}

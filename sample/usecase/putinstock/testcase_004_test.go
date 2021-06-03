package putinstock

import (
	"accounting/application/apperror"
	"accounting/domain/entity"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// WAVG
//
//input
//qty : 0
//totalPrice : 4200
//
//
// ERROR

type mockOutport4 struct {
	t         *testing.T
	nowDate   time.Time
	lastQty   int
	lastPrice float64
}

func (m mockOutport4) SaveStockPrice(ctx context.Context, obj *entity.StockPrice) error {

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
	assert.Equal(m.t, obj.Price, float64(2028))
	assert.Equal(m.t, obj.Date, m.nowDate)

	return nil
}

func (m mockOutport4) FindLastStockPrice(ctx context.Context, stockPriceID string) (*entity.StockPrice, error) {

	return &entity.StockPrice{
		ID:            "34",
		Date:          m.nowDate.Add(-1 * time.Hour),
		Quantity:      m.lastQty,
		Price:         m.lastPrice,
		InventoryCode: "X12",
		UpdatedByID:   "",
	}, nil

}

func (m mockOutport4) GenerateUUID(ctx context.Context) string {
	return "35"
}

func TestNormal4(t *testing.T) {

	now := time.Now()

	mockOutport4Obj := mockOutport4{
		t:         t,
		nowDate:   now,
		lastQty:   5,
		lastPrice: 2000,
	}

	ctx := context.Background()

	_, err := NewUsecase(&mockOutport4Obj).Execute(ctx, InportRequest{
		Date:          now,
		Quantity:      0,
		TotalPrice:    4200,
		InventoryCode: "X12",
		Method:        "WAVG",
	})

	assert.EqualError(t, err, apperror.QuantityMustNotZero.Error())

}

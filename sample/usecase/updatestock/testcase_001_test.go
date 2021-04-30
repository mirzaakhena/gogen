package updatestock

import (
	"accounting/domain/entity"
	"accounting/usecase/updatestock/mocks"
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCase007(t *testing.T) {

	{
		err := testedBody(t,
			"WAVG",
			2,
			1000,
			2,
			500,
			nil,
			[]*entity.StockPrice{
				{Quantity: 2, Price: 500},
			})

		assert.NoError(t, err)
	}

	{
		err := testedBody(t,
			"WAVG",
			2,
			1000,
			2,
			500,
			[]*entity.StockPrice{
				{Quantity: 3, Price: 500},
			},
			[]*entity.StockPrice{
				{Quantity: 5, Price: 500},
			})

		assert.NoError(t, err)
	}

	{
		err := testedBody(t,
			"WAVG",
			2,
			1000,
			2,
			500,
			[]*entity.StockPrice{
				{Quantity: 3, Price: 400},
			},
			[]*entity.StockPrice{
				{Quantity: 5, Price: 440},
			})

		assert.NoError(t, err)
	}

	{
		err := testedBody(t,
			"WAVG",
			-2,
			0,
			-2,
			500,
			nil,
			nil)

		assert.Error(t, err)
	}

	{
		err := testedBody(t,
			"WAVG",
			-2,
			0,
			-2,
			500,
			[]*entity.StockPrice{
				{Quantity: 1, Price: 500},
			},
			nil)

		assert.Error(t, err)
	}

	{
		err := testedBody(t,
			"FIFO",
			2,
			1000,
			2,
			500,
			nil,
			[]*entity.StockPrice{
				{Quantity: 2, Price: 500},
			})

		assert.Nil(t, err)
	}

	{
		err := testedBody(t,
			"FIFO",
			2,
			1000,
			2,
			500,
			[]*entity.StockPrice{
				{Quantity: 3, Price: 500},
			},
			[]*entity.StockPrice{
				{Quantity: 5, Price: 500},
			})

		assert.Nil(t, err)
	}

	{
		err := testedBody(t,
			"FIFO",
			2,
			1200,
			2,
			600,
			[]*entity.StockPrice{
				{Quantity: 3, Price: 500},
			},
			[]*entity.StockPrice{
				{Quantity: 3, Price: 500},
				{Quantity: 2, Price: 600},
			})

		assert.Nil(t, err)
	}

	{
		err := testedBody(t,
			"FIFO",
			2,
			1000,
			2,
			500,
			[]*entity.StockPrice{
				{Quantity: 3, Price: 500},
				{Quantity: 1, Price: 600},
			},
			[]*entity.StockPrice{
				{Quantity: 3, Price: 500},
				{Quantity: 1, Price: 600},
				{Quantity: 2, Price: 500},
			})

		assert.Nil(t, err)
	}

	{
		err := testedBody(t,
			"FIFO",
			2,
			1000,
			2,
			500,
			[]*entity.StockPrice{
				{Quantity: 1, Price: 600},
				{Quantity: 3, Price: 500},
			},
			[]*entity.StockPrice{
				{Quantity: 1, Price: 600},
				{Quantity: 5, Price: 500},
			})

		assert.Nil(t, err)
	}

	{
		err := testedBody(t,
			"FIFO",
			-2,
			0,
			-2,
			600,
			[]*entity.StockPrice{
				{Quantity: 5, Price: 600},
				{Quantity: 3, Price: 500},
			},
			[]*entity.StockPrice{
				{Quantity: 3, Price: 600},
				{Quantity: 3, Price: 500},
			})

		assert.Nil(t, err)
	}

}

// TestCase007 is for the case where
// We put the inventory in the second time
// but we already have the available stock before (with the same price)
//
// FIFO
// We have existing stock with Qty 3 and Price 500
// We want to pull stock with Qty 2
// the remaining stock must be Qty 1 with price 500
func testedBody(t *testing.T, calcMeth string, inputQty int, totalPrice float64, savedQty int, savedPrice float64, lastStockPrice []*entity.StockPrice, savedStockPrice []*entity.StockPrice) error {

	ctx := context.Background()

	mockOutport := mocks.Outport{}

	now := time.Now()

	mockOutport.
		On("SaveInventoryStock", ctx, &entity.InventoryBalance{
			ID:            "",
			Date:          now,
			ReferenceID:   "",
			Description:   "",
			InventoryCode: "12",
			Quantity:      savedQty,
			Price:         savedPrice,
			StockPrices:   savedStockPrice,
		}).
		Return(nil)

	// FindLastStockPrice return nothing here
	mockOutport.
		On("FindLastStockPrice", ctx, "12").
		Return(lastStockPrice, nil)

	mockOutport.
		On("FindOneInventory", ctx, "12").
		Return(&entity.Inventory{
			ID:                "12",
			CalculationMethod: calcMeth,
		}, nil)

	usecase := NewUsecase(&mockOutport)
	_, err := usecase.Execute(ctx, InportRequest{
		Date:          now,
		ReferenceID:   "",
		Description:   "",
		InventoryCode: "12",
		TotalPrice:    totalPrice,
		Quantity:      inputQty,
	})

	return err

}

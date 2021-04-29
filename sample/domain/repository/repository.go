package repository

import (
	"accounting/domain/entity"
	"accounting/domain/vo"
	"context"
)

type SaveJournalRepo interface {
	SaveJournal(ctx context.Context, obj *entity.Journal) error
}

type FindAllLastJournalBalanceRepo interface {

	// FindAllLastJournalBalance will return codes map of JournalBalance.
	// if the journal balance is not found then the value is 0.0
	// Error is happen when something goes wrong
	FindAllLastJournalBalance(ctx context.Context, businessID string, accountCodes []string) (map[string]float64, error)
}

type FindAllAccountSideByCodesRepo interface {

	// FindAllAccountSideByCodes will return all account side for respective account codes.
	// If account code not found then will return error
	FindAllAccountSideByCodes(ctx context.Context, businessID string, accountCode []string) (map[string]vo.AccountSide, error)
}

type SaveInventoryStockRepo interface {
	SaveInventoryStock(ctx context.Context, obj *entity.InventoryBalance) error
}

type FindLastStockPriceRepo interface {
	FindLastStockPrice(ctx context.Context, stockPriceID string) ([]*entity.StockPrice, error)
}

type FindOneInventoryRepo interface {
	FindOneInventory(ctx context.Context, inventoryID string) (*entity.Inventory, error)
}

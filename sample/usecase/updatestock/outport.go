package updatestock

import "accounting/domain/repository"

// Outport of UpdateStock
type Outport interface {
	repository.SaveInventoryStockRepo
	repository.FindLastStockPriceRepo
	repository.FindOneInventoryRepo
}

package pulloutstock

import (
	"accounting/domain/repository"
	"accounting/domain/service"
)

// Outport of PullOutStock
type Outport interface {
	repository.SaveStockPriceRepo
	repository.FindLastStockPriceRepo
	service.GenerateUUIDService
}

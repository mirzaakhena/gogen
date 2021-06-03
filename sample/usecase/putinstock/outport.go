package putinstock

import (
	"accounting/domain/repository"
	"accounting/domain/service"
)

// Outport of PutInStock
type Outport interface {
	repository.SaveStockPriceRepo
	repository.FindLastStockPriceRepo
	service.GenerateUUIDService
}

package createjournal

import (
	"accounting/domain/repository"
	"accounting/domain/service"
)

// Outport of CreateJournal
type Outport interface {
	service.GenerateUUIDService
	repository.SaveJournalRepo
	repository.FindLastJournalBalanceRepo
	repository.SaveJournalBalanceRepo
	repository.FindOneAccountByCodeRepo
}

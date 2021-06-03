package createjournal

import (
	"accounting/domain/repository"
	"accounting/domain/service"
)

// Outport of CreateJournal
type Outport interface {
	repository.TransactionRepo
	service.GenerateUUIDService
	repository.SaveJournalRepo
	repository.FindAllLastJournalBalanceRepo
	repository.FindAllAccountSideByCodesRepo
}

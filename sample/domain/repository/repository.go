package repository

import (
	"accounting/domain/entity"
	"context"
)

type SaveJournalRepo interface {
	SaveJournal(ctx context.Context, obj *entity.Journal) error
}

type FindLastJournalBalanceRepo interface {
	FindLastJournalBalance(ctx context.Context, bussinessID, accountCode string) (*entity.JournalBalance, error)
}

type SaveJournalBalanceRepo interface {
	SaveJournalBalance(ctx context.Context, obj *entity.JournalBalance) error
}

type FindOneAccountByCodeRepo interface {
	FindOneAccountByCode(ctx context.Context, accountID string) (*entity.Account, error)
}

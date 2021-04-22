package gateway

import (
	"accounting/domain/entity"
	"accounting/infrastructure/log2"
	"context"
)

type inmemoryGateway struct {
}

// NewInmemoryGateway ...
func NewInmemoryGateway() *inmemoryGateway {
	return &inmemoryGateway{}
}

func (r *inmemoryGateway) GenerateUUID(ctx context.Context) string {
	log2.Info(ctx, "called")

	return ""
}

func (r *inmemoryGateway) SaveJournal(ctx context.Context, obj *entity.Journal) error {
	log2.Info(ctx, "called")

	return nil
}

func (r *inmemoryGateway) FindLastJournalBalance(ctx context.Context, bussinessID, accountCode string) (*entity.JournalBalance, error) {
	log2.Info(ctx, "called")

	return nil, nil
}

func (r *inmemoryGateway) SaveJournalBalance(ctx context.Context, obj *entity.JournalBalance) error {
	log2.Info(ctx, "called")

	return nil
}

func (r *inmemoryGateway) FindOneAccountByCode(ctx context.Context, code string) (*entity.Account, error) {
	log2.Info(ctx, "called")

	return nil, nil
}

func (r *inmemoryGateway) BeginTransaction(ctx context.Context) context.Context {
	log2.Info(ctx, "called")

	return ctx
}

func (r *inmemoryGateway) CommitTransaction(ctx context.Context) {
	log2.Info(ctx, "called")

	return
}

func (r *inmemoryGateway) RollbackTransaction(ctx context.Context) {
	log2.Info(ctx, "called")

	return
}

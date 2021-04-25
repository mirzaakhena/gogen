package gateway

import (
	"accounting/domain/entity"
	"accounting/domain/vo"
	"accounting/infrastructure/log"
	"context"
)

type inmemoryGateway struct {
}

// NewInmemoryGateway ...
func NewInmemoryGateway() *inmemoryGateway {
	return &inmemoryGateway{}
}

func (r *inmemoryGateway) GenerateUUID(ctx context.Context) string {
	log.Info(ctx, "called")

	return ""
}

func (r *inmemoryGateway) SaveJournal(ctx context.Context, obj *entity.Journal) error {
	log.Info(ctx, "called")

	return nil
}

func (r *inmemoryGateway) BeginTransaction(ctx context.Context) context.Context {
	log.Info(ctx, "called")

	return ctx
}

func (r *inmemoryGateway) CommitTransaction(ctx context.Context) {
	log.Info(ctx, "called")

	return
}

func (r *inmemoryGateway) RollbackTransaction(ctx context.Context) {
	log.Info(ctx, "called")

	return
}

func (r *inmemoryGateway) FindAllLastJournalBalance(ctx context.Context, businessID string, accountCodes []string) (map[string]float64, error) {
	log.Info(ctx, "called")

	return nil, nil
}

func (r *inmemoryGateway) FindAllAccountSideByCodes(ctx context.Context, businessID string, accountCode []string) (map[string]vo.AccountSide, error) {
	log.Info(ctx, "called")

	return nil, nil
}

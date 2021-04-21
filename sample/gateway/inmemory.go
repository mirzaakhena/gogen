package gateway

import (
	"accounting/domain/entity"
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
	log.InfoRequest(ctx, "called")

	return ""
}

func (r *inmemoryGateway) SaveJournal(ctx context.Context, obj *entity.Journal) error {
	log.InfoRequest(ctx, "called")

	return nil
}

func (r *inmemoryGateway) FindLastJournalBalance(ctx context.Context, bussinessID, accountCode string) (*entity.JournalBalance, error) {
	log.InfoRequest(ctx, "called")

	return nil, nil
}

func (r *inmemoryGateway) SaveJournalBalance(ctx context.Context, obj *entity.JournalBalance) error {
	log.InfoRequest(ctx, "called")

	return nil
}

func (r *inmemoryGateway) FindOneAccountByCode(ctx context.Context, code string) (*entity.Account, error) {
	log.InfoRequest(ctx, "called")

	return nil, nil
}
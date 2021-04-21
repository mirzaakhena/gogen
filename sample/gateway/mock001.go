package gateway

import (
	"accounting/domain/entity"
	"accounting/infrastructure/log"
	"context"
)

type mock001Gateway struct {
}

// NewMock001Gateway ...
func NewMock001Gateway() *mock001Gateway {
	return &mock001Gateway{}
}

func (r *mock001Gateway) GenerateUUID(ctx context.Context) string {
	log.InfoRequest(ctx, "called")

	return "345"
}

func (r *mock001Gateway) SaveJournal(ctx context.Context, obj *entity.Journal) error {
	log.InfoRequest(ctx, "called")

	return nil
}

func (r *mock001Gateway) FindLastJournalBalance(ctx context.Context, bussinessID, accountCode string) (*entity.JournalBalance, error) {
	log.InfoRequest(ctx, "called")

	return nil, nil
}

func (r *mock001Gateway) SaveJournalBalance(ctx context.Context, obj *entity.JournalBalance) error {
	log.InfoRequest(ctx, "called")

	return nil
}

func (r *mock001Gateway) FindOneAccountByCode(ctx context.Context, code string) (*entity.Account, error) {
	log.InfoRequest(ctx, "called")

	if code == "100" {
		return &entity.Account{
			ID:   "12",
			Side: "ACTIVA",
		}, nil

	} else if code == "101" {
		return &entity.Account{
			ID:   "13",
			Side: "PASSIVA",
		}, nil
	}

	return nil, nil
}

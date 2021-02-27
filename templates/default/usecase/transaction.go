package usecase

import (
	"context"
)

type Transaction interface {
	BeginTransaction(ctx context.Context) (context.Context, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

func ExecuteTransaction(ctx context.Context, trx Transaction, trxFunc func(ctx context.Context) error) error {
	dbCtx, err := trx.BeginTransaction(ctx)
	if err != nil {
		return err
	}

	defer func() {
		if p := recover(); p != nil {
			err = trx.Rollback(dbCtx)
			panic(p)

		} else if err != nil {
			err = trx.Rollback(dbCtx)

		} else {
			err = trx.Commit(dbCtx)

		}
	}()

	err = trxFunc(dbCtx)
	return err
}

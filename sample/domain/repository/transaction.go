package repository

import "context"

// TransactionRepo used for common transaction handling
// all the context must use the same database session.
type TransactionRepo interface {

	// BeginTransaction is wrapping the transaction database object to context and return it
	BeginTransaction(ctx context.Context) context.Context

	// CommitTransaction will save the changes
	CommitTransaction(ctx context.Context)

	// RollbackTransaction will undo the changes
	RollbackTransaction(ctx context.Context)
}

// ExecuteTransaction is helper function that simplify the transaction execution handling
func ExecuteTransaction(ctx context.Context, trx TransactionRepo, trxFunc func(ctx context.Context) (err error)) (err error) {
	dbCtx := trx.BeginTransaction(ctx)

	defer func() {
		if p := recover(); p != nil {
			trx.RollbackTransaction(dbCtx)

		} else if err != nil {
			trx.RollbackTransaction(dbCtx)

		} else {
			trx.CommitTransaction(dbCtx)

		}
	}()

	err = trxFunc(dbCtx)

	return
}

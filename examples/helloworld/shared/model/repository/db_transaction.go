package repository

import "context"

// WithoutTransactionDB used to get database object from any database implementation.
// For consistency reason both WithTransactionDB and WithoutTransactionDB will seek database object under the context params
type WithoutTransactionDB interface {
	GetDatabase(ctx context.Context) (context.Context, error)
	Close(ctx context.Context) error
}

// WithTransactionDB used for common transaction handling
// all the context must use the same database session.
type WithTransactionDB interface {
	BeginTransaction(ctx context.Context) (context.Context, error)
	CommitTransaction(ctx context.Context) error
	RollbackTransaction(ctx context.Context) error
}

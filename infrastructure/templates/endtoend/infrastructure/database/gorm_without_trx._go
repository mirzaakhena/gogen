package database

import (
	"context"
	"gorm.io/gorm"
)

type GormWithoutTransactionImpl struct {
	db *gorm.DB
}

func NewGormWithoutTransactionImpl(db *gorm.DB) *GormWithoutTransactionImpl {
	return &GormWithoutTransactionImpl{
		db: db,
	}
}

func (r *GormWithoutTransactionImpl) GetDatabase(ctx context.Context) (context.Context, error) {
	trxCtx := context.WithValue(ctx, ContextDBValue, r.db)
	return trxCtx, nil
}

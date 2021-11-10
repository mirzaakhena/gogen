package database

import (
	"context"
	"fmt"
	"gorm.io/gorm"
)

type contextDBType string

var ContextDBValue contextDBType = "DB"

// ExtractDB is used by other repo to extract the databasex from context
func ExtractDB(ctx context.Context) (*gorm.DB, error) {

	db, ok := ctx.Value(ContextDBValue).(*gorm.DB)
	if !ok {
		return nil, fmt.Errorf("database is not found in context")
	}

	return db, nil
}

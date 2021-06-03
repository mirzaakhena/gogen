package service

import (
	"context"
)

type GenerateUUIDService interface {
	GenerateUUID(ctx context.Context) string
}

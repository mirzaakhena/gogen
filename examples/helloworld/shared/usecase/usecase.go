package usecase

import "context"

type Inport[REQUEST, RESPONSE any] interface {
	Execute(ctx context.Context, req REQUEST) (*RESPONSE, error)
}

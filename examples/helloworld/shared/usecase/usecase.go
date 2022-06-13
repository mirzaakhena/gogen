package usecase

type Inport[CONTEXT, REQUEST, RESPONSE any] interface {
	Execute(ctx CONTEXT, req REQUEST) (*RESPONSE, error)
}

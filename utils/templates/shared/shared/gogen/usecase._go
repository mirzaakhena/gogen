package gogen

import (
	"context"
	"fmt"
	"os"
)

type Inport[REQUEST, RESPONSE any] interface {
	Execute(ctx context.Context, req REQUEST) (*RESPONSE, error)
}

func GetInport[Req, Res any](usecase any, err error) Inport[Req, Res] {

	if err != nil {
		fmt.Printf("\n\n%s...\n\n", err.Error())
		os.Exit(0)
	}

	inport, ok := usecase.(Inport[Req, Res])
	if !ok {
		fmt.Printf("unable to cast to Inport\n")
		os.Exit(0)
	}
	return inport
}

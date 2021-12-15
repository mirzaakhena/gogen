package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/genusecase"
)

// genTestHandler ...
func (r *Controller) genUsecaseHandler(inputPort genusecase.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("\n" +
				"   # Create a new usecase\n" +
				"   gogen usecase CreateOrder\n" +
				"     'CreateOrder' is an usecase name\n" +
				"\n")

			return err
		}

		usecaseName := commands[0]

		var req genusecase.InportRequest
		req.UsecaseName = usecaseName

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

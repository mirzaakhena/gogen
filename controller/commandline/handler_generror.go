package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/generror"
)

// genErrorHandler ...
func (r *Controller) genErrorHandler(inputPort generror.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("\n" +
				"   # Create an error enum\n" +
				"   gogen error SomethingGoesWrongError\n" +
				"     'SomethingGoesWrongError' is an error constant name\n" +
				"\n")

			return err
		}

		var req generror.InportRequest
		req.ErrorName = commands[0]

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

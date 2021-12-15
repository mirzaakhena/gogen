package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/gentest"
)

// genTestHandler ...
func (r *Controller) genTestHandler(inputPort gentest.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 2 {
			err := fmt.Errorf("\n" +
				"   # Create a test case file for current usecase\n" +
				"   gogen test normal CreateOrder\n" +
				"     'normal'      is a test case name\n" +
				"     'CreateOrder' is an usecase name\n" +
				"\n")

			return err
		}

		var req gentest.InportRequest
		req.TestName = commands[0]
		req.UsecaseName = commands[1]

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

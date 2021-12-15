package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/genendtoend"
)

// genEndToEndHandler ...
func (r *Controller) genEndToEndHandler(inputPort genendtoend.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("\n" +
				"   # Create a complete Usecase Gateway Controller and registry bootstrap code\n" +
				"   gogen endtoend CreateOrder Order\n" +
				"     'CreateOrder' is an usecase name\n" +
				"     'Order' is an entity name\n" +
				"\n")
			return err
		}

		var req genendtoend.InportRequest
		req.UsecaseName = commands[0]
		req.EntityName = commands[1]

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

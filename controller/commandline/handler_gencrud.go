package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/gencrud"
)

// genCrudHandler ...
func (r *Controller) genCrudHandler(inputPort gencrud.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("invalid gogen crud command format. Try this `gogen crud EntityName`")
			return err
		}

		var req gencrud.InportRequest
		req.EntityName = commands[0]

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

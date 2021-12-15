package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/genwebapp"
)

// genWebappHandler ...
func (r *Controller) genWebappHandler(inputPort genwebapp.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("\n" +
				"   # Create a complete CRUD webapp ui for specific entity\n" +
				"   gogen webapp Product\n" +
				"     'Product' is an entity name\n" +
				"\n")
			return err
		}

		var req genwebapp.InportRequest
		req.EntityName = commands[0]

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

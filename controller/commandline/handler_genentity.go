package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/genentity"
)

// genEntityHandler ...
func (r *Controller) genEntityHandler(inputPort genentity.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("\n" +
				"   # Create an entity\n" +
				"   gogen entity Product\n" +
				"     'Product' is an entity name\n" +
				"\n")
			return err
		}

		var req genentity.InportRequest
		req.EntityName = commands[0]

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

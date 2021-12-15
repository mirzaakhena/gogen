package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/genrepository"
)

// genRepositoryHandler ...
func (r *Controller) genRepositoryHandler(inputPort genrepository.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 2 {
			err := fmt.Errorf("\n" +
				"   # Create a repository and inject the template code into interactor file with '//!' flag\n" +
				"   gogen repository SaveOrder Order CreateOrder\n" +
				"     'SaveOrder'   is a repository func name\n" +
				"     'Order'       is an entity name\n" +
				"     'CreateOrder' is an usecase name\n" +
				"\n" +
				"   # Create a repository without inject the template code into usecase\n" +
				"   gogen repository SaveOrder Order\n" +
				"     'SaveOrder' is a repository func name\n" +
				"     'Order'     is an entity name\n" +
				"\n")

			return err
		}

		var req genrepository.InportRequest
		req.RepositoryName = commands[0]
		req.EntityName = commands[1]

		if len(commands) >= 3 {
			req.UsecaseName = commands[2]
		}

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

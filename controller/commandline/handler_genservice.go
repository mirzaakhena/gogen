package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/genservice"
)

// genServiceHandler ...
func (r *Controller) genServiceHandler(inputPort genservice.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("   # Create a service and inject the template code into interactor file with '//!' flag\n" +
				"   gogen service PublishMessage CreateOrder\n" +
				"     'PublishMessage' is a service func name\n" +
				"     'CreateOrder'    is an usecase name\n" +
				"\n" +
				"   # Create a service without inject the template code into usecase\n" +
				"   gogen service PublishMessage\n" +
				"     'PublishMessage' is a service func name\n" +
				"\n")
			return err
		}

		var req genservice.InportRequest
		req.ServiceName = commands[0]

		if len(commands) >= 2 {
			req.UsecaseName = commands[1]
		}

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

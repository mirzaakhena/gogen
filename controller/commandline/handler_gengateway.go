package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/gengateway"
)

// genGatewayHandler ...
func (r *Controller) genGatewayHandler(inputPort gengateway.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("\n" +
				"   # Create a gateway for specific usecase\n" +
				"   gogen gateway inmemory CreateOrder\n" +
				"     'inmemory'    is a gateway name\n" +
				"     'CreateOrder' is an usecase name\n" +
				"\n" +
				"   # Create a gateway for all usecases\n" +
				"   gogen gateway inmemory\n" +
				"     'inmemory' is a gateway name\n" +
				"\n")

			return err
		}

		var req gengateway.InportRequest
		req.GatewayName = commands[0]

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

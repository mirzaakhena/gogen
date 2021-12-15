package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/gencontroller"
)

// genControllerHandler ...
func (r *Controller) genControllerHandler(inputPort gencontroller.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 2 {

			err := fmt.Errorf("\n" +
				"   # Create a controller with defined web framework or other handler\n" +
				"   gogen controller restapi CreateOrder gin\n" +
				"     'restapi'     is a gateway name\n" +
				"     'CreateOrder' is an usecase name\n" +
				"     'gin'         is a sample webframewrok. You may try the other one like: nethttp, echo, and gorilla\n" +
				"\n" +
				"   # Create a controller with gin as default web framework\n" +
				"   gogen controller restapi CreateOrder\n" +
				"     'restapi'     is a gateway name\n" +
				"     'CreateOrder' is an usecase name\n" +
				"\n")

			return err
		}

		var req gencontroller.InportRequest

		req.ControllerName = commands[0]

		req.UsecaseName = commands[1]

		if len(commands) >= 3 {
			req.DriverName = commands[2]
		} else {
			req.DriverName = "gin"
		}

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

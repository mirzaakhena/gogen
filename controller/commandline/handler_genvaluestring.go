package commandline

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/usecase/genvaluestring"
)

// genValueStringHandler ...
func (r *Controller) genValueStringHandler(inputPort genvaluestring.Inport) func(...string) error {

	return func(commands ...string) error {

		ctx := context.Background()

		if len(commands) < 1 {
			err := fmt.Errorf("\n" +
				"   # Create a valueobject with simple string type\n" +
				"   gogen valuestring OrderID\n" +
				"     'OrderID' is an valueobject name\n" +
				"\n")
			return err
		}

		var req genvaluestring.InportRequest
		req.ValueStringName = commands[0]

		_, err := inputPort.Execute(ctx, req)
		if err != nil {
			return err
		}

		return nil

	}
}

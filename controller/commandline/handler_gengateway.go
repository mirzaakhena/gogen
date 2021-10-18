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

    if len(commands) < 2 {
      err := fmt.Errorf("invalid gogen gateway format. Try this `gogen gateway GatewayName UsecaseName`")
      return err
    }

    var req gengateway.InportRequest
    req.GatewayName = commands[0]
    req.UsecaseName = commands[1]

    _, err := inputPort.Execute(ctx, req)
    if err != nil {
      return err
    }

    return nil

  }
}

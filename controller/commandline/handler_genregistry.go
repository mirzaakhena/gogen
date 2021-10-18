package commandline

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/usecase/genregistry"
)

// genRegistryHandler ...
func (r *Controller) genRegistryHandler(inputPort genregistry.Inport) func(...string) error {

  return func(commands ...string) error {

    ctx := context.Background()

    if len(commands) < 1 {
      err := fmt.Errorf("invalid gogen registry command format. Try this `gogen registry RegistryName`")
      return err
    }

    var req genregistry.InportRequest
    req.RegistryName = commands[0]

    if len(commands) >= 2 {
      req.ControllerName = commands[1]
    }

    if len(commands) >= 3 {
      req.GatewayName = commands[2]
    }

    _, err := inputPort.Execute(ctx, req)
    if err != nil {
      return err
    }

    return nil

  }
}

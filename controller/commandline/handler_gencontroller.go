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
      err := fmt.Errorf("invalid gogen controller command format. Try this `gogen controller ControllerName UsecaseName`")
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

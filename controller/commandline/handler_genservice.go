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
      err := fmt.Errorf("invalid gogen service command format. Try this `gogen service ServiceName [UsecaseName]`")
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

package commandline

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/usecase/gentest"
)

// genTestHandler ...
func (r *Controller) genTestHandler(inputPort gentest.Inport) func(...string) error {

  return func(commands ...string) error {

    ctx := context.Background()

    if len(commands) < 2 {
      err := fmt.Errorf("invalid gogen test command format. Try this `gogen test normal UsecaseName`")
      return err
    }

    var req gentest.InportRequest
    req.TestName = commands[0]
    req.UsecaseName = commands[1]

    _, err := inputPort.Execute(ctx, req)
    if err != nil {
      return err
    }

    return nil

  }
}

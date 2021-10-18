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
      err := fmt.Errorf("invalid gogen valuestring command format. Try this `gogen valuestring VOName`")
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

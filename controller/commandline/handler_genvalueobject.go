package commandline

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/usecase/genvalueobject"
)

// genValueObjectHandler ...
func (r *Controller) genValueObjectHandler(inputPort genvalueobject.Inport) func(...string) error {

  return func(commands ...string) error {

    ctx := context.Background()

    if len(commands) < 2 {
      err := fmt.Errorf("invalid gogen valueobject command format. Try this `gogen valueobject VOName Field1 Field2 ...`")
      return err
    }

    var req genvalueobject.InportRequest
    req.ValueObjectName = commands[0]
    req.FieldNames = commands[1:]

    _, err := inputPort.Execute(ctx, req)
    if err != nil {
      return err
    }

    return nil

  }
}

package commandline

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/usecase/genrepository"
)

// genRepositoryHandler ...
func (r *Controller) genRepositoryHandler(inputPort genrepository.Inport) func(...string) error {

  return func(commands ...string) error {

    ctx := context.Background()

    if len(commands) < 2 {
      err := fmt.Errorf("invalid gogen repository command format. Try this `gogen repository RepoName EntityName [UsecaseName]`")
      return err
    }

    var req genrepository.InportRequest
    req.RepositoryName = commands[0]
    req.EntityName = commands[1]

    if len(commands) >= 3 {
      req.UsecaseName = commands[2]
    }

    _, err := inputPort.Execute(ctx, req)
    if err != nil {
      return err
    }

    return nil

  }
}

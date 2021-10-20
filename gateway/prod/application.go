package prod

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/infrastructure/templates"
)

func (r *prodGateway) GetApplicationFileTemplate(ctx context.Context, driverName string) string {
  return templates.ReadFile(fmt.Sprintf("application/registry/~driver_%s._go", driverName))
}

// GetErrorLineTemplate ...
func (r prodGateway) GetErrorLineTemplate(ctx context.Context) string {
	return templates.ReadFile("application/apperror/~inject._go")
}

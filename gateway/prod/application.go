package prod

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/infrastructure/templates"
)

func (r *prodGateway) GetApplicationFileTemplate(ctx context.Context) string {
	return templates.ReadFile("application/application._go")
}

func (r *prodGateway) GetRegistryFileTemplate(ctx context.Context, driverName string) string {
	return templates.ReadFile(fmt.Sprintf("application/registry/~driver_%s._go", driverName))
}

// GetErrorLineTemplate ...
func (r prodGateway) GetErrorLineTemplate(ctx context.Context) string {
	return templates.ReadFile("domain/domerror/~inject._go")
}

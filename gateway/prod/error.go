package prod

import (
	"context"
	"github.com/mirzaakhena/gogen/infrastructure/templates"
)

type errorGateway struct {
}

func (r errorGateway) GetErrorEnumTemplate(ctx context.Context) string {
	return templates.ReadFile("domain/domerror/error_enum._go")
}

func (r errorGateway) GetErrorFuncTemplate(ctx context.Context) string {
	return templates.ReadFile("domain/domerror/error_func._go")
}

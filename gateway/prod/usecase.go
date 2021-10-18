package prod

import (
	"context"
  "github.com/mirzaakhena/gogen/infrastructure/templates"
)

// GetInportTemplate ...
func (r *prodGateway) GetInportTemplate(ctx context.Context) string {
	return templates.ReadFile("usecase/usecase_inport._go")
}

// GetOutportTemplate ...
func (r *prodGateway) GetOutportTemplate(ctx context.Context) string {
	return templates.ReadFile("usecase/usecase_outport._go")
}

// GetInteractorTemplate ...
func (r *prodGateway) GetInteractorTemplate(ctx context.Context) string {
	return templates.ReadFile("usecase/usecase_interactor._go")
}

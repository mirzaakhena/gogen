package prod

import (
	"context"
	"github.com/mirzaakhena/gogen/infrastructure/templates"
)

//// GetRepositoryTemplate ...
//func (r *prodGateway) GetRepositoryTemplate(ctx context.Context) string {
//  return templates.ReadFile("model/repository/repository._go")
//}

// GetServiceInterfaceTemplate ...
func (r *prodGateway) GetServiceInterfaceTemplate(ctx context.Context) string {
	return templates.ReadFile("model/service/~service_interface._go")
}

// GetServiceInjectTemplate ...
func (r *prodGateway) GetServiceInjectTemplate(ctx context.Context) string {
	return templates.ReadFile("model/service/~service_inject._go")
}

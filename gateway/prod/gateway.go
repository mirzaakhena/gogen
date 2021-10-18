package prod

import (
  "context"
  "github.com/mirzaakhena/gogen/infrastructure/templates"
)

// GetGatewayTemplate ...
func (r *prodGateway) GetGatewayTemplate(ctx context.Context) string {
  return templates.ReadFile("gateway/gorm/impl._go")
}

// GetGatewayMethodTemplate ...
func (r *prodGateway) GetGatewayMethodTemplate(ctx context.Context) string {
  return templates.ReadFile("gateway/${gatewayname}/~inject._go")
}

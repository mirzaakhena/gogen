package prod

import (
  "context"
  "github.com/mirzaakhena/gogen/infrastructure/templates"
)

// GetValueObjectTemplate ...
func (r *prodGateway) GetValueObjectTemplate(ctx context.Context) string {
  return templates.ReadFile("domain/vo/domain_vo_valueobject._go")
}

// GetValueStringTemplate ...
func (r *prodGateway) GetValueStringTemplate(ctx context.Context) string {
  return templates.ReadFile("domain/vo/domain_vo_valuestring._go")
}

package prod

import (
	"context"
  "github.com/mirzaakhena/gogen/infrastructure/templates"
)

// GetTestTemplate ...
func (r *prodGateway) GetTestTemplate(ctx context.Context) string {
	return templates.ReadFile("usecase/usecase_test._go")
}

// GetTestMethodTemplate ...
func (r *prodGateway) GetTestMethodTemplate(ctx context.Context) string {
	return templates.ReadFile("test/usecase/${usecasename}/~inject._go")
}

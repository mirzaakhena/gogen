package generror

import (
	"context"
  "github.com/mirzaakhena/gogen/domain/service"
)

// Outport of GenError
type Outport interface {
	service.PrintTemplateService
	service.ReformatService

	GetErrorLineTemplate(ctx context.Context) string
}

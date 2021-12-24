package generror

import (
	"context"
	"github.com/mirzaakhena/gogen/model/service"
)

// Outport of GenError
type Outport interface {
	service.PrintTemplateService
	service.ReformatService

	GetErrorLineTemplate(ctx context.Context) string
}

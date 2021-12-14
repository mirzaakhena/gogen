package gencrud

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/service"
)

// Outport of usecase
type Outport interface {
	service.GetPackagePathService
	GetMainFileForCrudTemplate(ctx context.Context) string
	service.WriteFileIfNotExistService
}

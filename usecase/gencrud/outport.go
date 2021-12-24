package gencrud

import (
	"context"
	"github.com/mirzaakhena/gogen/model/service"
)

// Outport of usecase
type Outport interface {
	service.GetPackagePathService
	GetMainFileForCrudTemplate(ctx context.Context) string
	service.WriteFileIfNotExistService
}

package genendtoend

import (
	"context"
	"github.com/mirzaakhena/gogen/model/service"
)

// Outport of usecase
type Outport interface {
	service.GetPackagePathService
	service.WriteFileIfNotExistService

	GetMainFileForE2ETemplate(ctx context.Context) string
}

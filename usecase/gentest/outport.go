package gentest

import (
	"context"
  "github.com/mirzaakhena/gogen/domain/service"
)

// Outport of GenTest
type Outport interface {
	service.LogActionInterface
	service.IsFileExistService
	service.WriteFileService
	service.ReformatService
	service.GetPackagePathService
	service.PrintTemplateService

	GetTestTemplate(ctx context.Context) string
	GetTestMethodTemplate(ctx context.Context) string
}

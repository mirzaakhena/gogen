package genregistry

import (
  "context"
  "github.com/mirzaakhena/gogen/domain/service"
)

// Outport of GenRegistry
type Outport interface {
  service.GetPackagePathService
  service.WriteFileIfNotExistService
  service.ReformatService
  GetMainFileTemplate(ctx context.Context) string
  GetServerFileTemplate(ctx context.Context, driverName string) string
  GetApplicationFileTemplate(ctx context.Context, driverName string) string
}

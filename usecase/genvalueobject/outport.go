package genvalueobject

import (
  "context"
  "github.com/mirzaakhena/gogen/domain/service"
)

// Outport of GenValueObject
type Outport interface {
  service.CreateFolderIfNotExistService
  service.WriteFileIfNotExistService

  GetValueObjectTemplate(ctx context.Context) string
}

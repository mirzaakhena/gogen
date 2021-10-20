package genvaluestring

import (
  "context"
  "github.com/mirzaakhena/gogen/domain/service"
)

// Outport of GenValueString
type Outport interface {
  service.CreateFolderIfNotExistService
  service.WriteFileIfNotExistService

  GetValueStringTemplate(ctx context.Context) string
}
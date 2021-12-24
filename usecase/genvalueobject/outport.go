package genvalueobject

import (
	"context"
	"github.com/mirzaakhena/gogen/model/service"
)

// Outport of GenValueObject
type Outport interface {
	service.CreateFolderIfNotExistService
	service.WriteFileIfNotExistService

	GetValueObjectTemplate(ctx context.Context) string
}

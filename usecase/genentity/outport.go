package genentity

import (
	"github.com/mirzaakhena/gogen/model/service"
)

// Outport of GenEntity
type Outport interface {
	service.CreateFolderIfNotExistService
	service.WriteFileIfNotExistService
	service.ReformatService
}

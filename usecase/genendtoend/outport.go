package genendtoend

import (
	"github.com/mirzaakhena/gogen/domain/service"
)

// Outport of usecase
type Outport interface {
	service.GetPackagePathService
	service.WriteFileIfNotExistService
}

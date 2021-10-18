package gengateway

import (
	"context"
  "github.com/mirzaakhena/gogen/domain/service"
)

// Outport of GenGateway
type Outport interface {
	// service.ApplicationActionInterface
	service.LogActionInterface
	service.GetPackagePathService
	service.IsFileExistService
	service.WriteFileService
	service.PrintTemplateService
	service.ReformatService

	GetGatewayMethodTemplate(ctx context.Context) string
	GetGatewayTemplate(ctx context.Context) string
}

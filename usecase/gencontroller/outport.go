package gencontroller

import (
	"context"
	"github.com/mirzaakhena/gogen/model/service"
)

// Outport of GenController
type Outport interface {
	service.ErrorActionInterface
	service.LogActionInterface

	service.IsFileExistService
	service.WriteFileService
	service.ReformatService
	service.GetPackagePathService
	service.PrintTemplateService

	//GetResponseTemplate(ctx context.Context) string
	//GetInterceptorTemplate(ctx context.Context, driverName string) string
	//GetRouterTemplate(ctx context.Context, driverName string) string

	GetHandlerTemplate(ctx context.Context, driverName string) string
	GetRouterInportTemplate(ctx context.Context, driverName string) string
	GetRouterRegisterTemplate(ctx context.Context, driverName string) string
}

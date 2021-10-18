package gencontroller

import (
  "context"
  "github.com/mirzaakhena/gogen/domain/service"
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

  GetControllerTemplate(ctx context.Context) string
  GetResponseTemplate(ctx context.Context) string
  GetInterceptorTemplate(ctx context.Context, framework string) string
  GetRouterTemplate(ctx context.Context, framework string) string
  GetHandlerTemplate(ctx context.Context, framework string) string
  GetRouterInportTemplate(ctx context.Context) string
  GetRouterRegisterTemplate(ctx context.Context) string
}

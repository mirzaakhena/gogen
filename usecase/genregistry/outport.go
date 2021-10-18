package genregistry

import (
  "context"
  "github.com/mirzaakhena/gogen/domain/service"
)

// Outport of GenRegistry
type Outport interface {
  service.GetPackagePathService
  GetMainFileTemplate(ctx context.Context) string
  service.WriteFileIfNotExistService
  //service.CreateFolderIfNotExistService
  //service.ReformatService

  //GetApplicationTemplate(ctx context.Context) string
  //GetRegistryTemplate(ctx context.Context) string
  //
  //FindObjController(ctx context.Context, controllerName string) (*entity.ObjController, error)
  //FindAllObjGateway(ctx context.Context) ([]*entity.ObjGateway, error)
  //FindObjGateway(ctx context.Context, gatewayName string) (*entity.ObjGateway, error)
  //FindAllObjUsecases(ctx context.Context, objController *entity.ObjController) ([]*entity.ObjUsecase, error)
}

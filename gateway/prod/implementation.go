package prod

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/domain/entity"
  "github.com/mirzaakhena/gogen/infrastructure/templates"
)

type prodGateway struct {
  *basicUtilityGateway
  *errorGateway
}

func (r *prodGateway) GetServerFileTemplate(ctx context.Context, driverName string) string {
  return templates.ReadFile(fmt.Sprintf("infrastructure/server/~http_server_%s._go", driverName))
}

func (r *prodGateway) GetRegistryTemplate(ctx context.Context) string {
  return templates.RegistryGingonicFile
}

func (r *prodGateway) GetMainFileTemplate(ctx context.Context) string {
  return templates.ReadFile("main._go")
}

func (r *prodGateway) FindObjController(ctx context.Context, controllerName string) (*entity.ObjController, error) {
  panic("implement me")
}

func (r *prodGateway) FindAllObjGateway(ctx context.Context) ([]*entity.ObjGateway, error) {
  panic("implement me")
}

func (r *prodGateway) FindObjGateway(ctx context.Context, gatewayName string) (*entity.ObjGateway, error) {
  panic("implement me")
}

func (r *prodGateway) FindAllObjUsecases(ctx context.Context, objController *entity.ObjController) ([]*entity.ObjUsecase, error) {
  panic("implement me")
}

// NewProdGateway ...
func NewProdGateway() *prodGateway {
  return &prodGateway{
    basicUtilityGateway: &basicUtilityGateway{},
    errorGateway:        &errorGateway{},
  }
}

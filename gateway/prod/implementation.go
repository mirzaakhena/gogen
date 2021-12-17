package prod

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/infrastructure/templates"
)

type prodGateway struct {
	*basicUtilityGateway
	*errorGateway
}

func (r *prodGateway) GetServerFileTemplate(ctx context.Context, driverName string) string {
	return templates.ReadFile(fmt.Sprintf("infrastructure/server/~http_server_%s._go", driverName))
}

//func (r *prodGateway) GetRegistryTemplate(ctx context.Context) string {
//	return templates.RegistryGingonicFile
//}

func (r *prodGateway) GetMainFileTemplate(ctx context.Context) string {
	return templates.ReadFile("main._go")
}

func (r *prodGateway) GetMainFileForCrudTemplate(ctx context.Context) string {
	return templates.ReadFileCrud("main._go")
}

func (r *prodGateway) GetMainFileForE2ETemplate(ctx context.Context) string {
	return templates.ReadFileE2E("main._go")
}

// NewProdGateway ...
func NewProdGateway() *prodGateway {
	return &prodGateway{
		basicUtilityGateway: &basicUtilityGateway{},
		errorGateway:        &errorGateway{},
	}
}

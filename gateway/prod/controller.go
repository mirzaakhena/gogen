package prod

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/infrastructure/templates"
)

//func (r *prodGateway) GetResponseTemplate(ctx context.Context) string {
//  return templates.ControllerGinGonicResponseFile
//}
//
//func (r *prodGateway) GetInterceptorTemplate(ctx context.Context, framework string) string {
//  return templates.ControllerGinGonicInterceptorFile
//}
//
//func (r *prodGateway) GetRouterTemplate(ctx context.Context, framework string) string {
//  return templates.ControllerGinGonicRouterStructFile
//}

func (r *prodGateway) GetHandlerTemplate(ctx context.Context, framework string) string {
	return templates.ControllerGinGonicHandlerFuncFile
}

func (r *prodGateway) GetRouterRegisterTemplate(ctx context.Context, driverName string) string {
	path := fmt.Sprintf("controllers/%s/controller/${controllername}/~inject-router-register._go", driverName)
	return templates.ReadFile(path)
}

func (r *prodGateway) GetRouterInportTemplate(ctx context.Context, driverName string) string {
	path := fmt.Sprintf("controllers/%s/controller/${controllername}/~inject-router-inport._go", driverName)
	return templates.ReadFile(path)
}

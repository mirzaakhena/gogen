package gencontroller

import (
	"context"
	"fmt"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
)

//go:generate mockery --name Outport -output mocks/

type genControllerInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenController
func NewUsecase(outputPort Outport) Inport {
	return &genControllerInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenController
func (r *genControllerInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	packagePath := r.outport.GetPackagePath(ctx)

	err := service.CreateEverythingExactly("default/", "application/apperror", map[string]string{}, struct{ PackagePath string }{PackagePath: packagePath})
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("default/", "infrastructure/log", map[string]string{}, struct{ PackagePath string }{PackagePath: packagePath})
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("default/", "infrastructure/util", map[string]string{}, struct{ PackagePath string }{PackagePath: packagePath})
	if err != nil {
		return nil, err
	}

	objUsecase, err := entity.NewObjUsecase(req.UsecaseName)
	if err != nil {
		return nil, err
	}

	objCtrl, err := entity.NewObjController(req.ControllerName)
	if err != nil {
		return nil, err
	}

	if req.DriverName == "" {
		req.DriverName = "gin"
	}

	skippedFolder := fmt.Sprintf("default/controllers/%s/", req.DriverName)
	err = service.CreateEverythingExactly(skippedFolder, "controller", map[string]string{
		"controllername": objCtrl.ControllerName.LowerCase(),
		"usecasename":    objUsecase.UsecaseName.LowerCase(),
	}, objCtrl.GetData(packagePath, *objUsecase))
	if err != nil {
		return nil, err
	}

	objDataCtrl := objCtrl.GetData(packagePath, *objUsecase)

	//inject inport to struct
	//type Controller struct {
	//  Router            gin.IRouter
	//  CreateOrderInport createorder.Inport <----- here
	//}
	{
		templateCode := r.outport.GetRouterInportTemplate(ctx, req.DriverName)

		templateWithData, err := r.outport.PrintTemplate(ctx, templateCode, objDataCtrl)
		if err != nil {
			return nil, err
		}

		bytes, err := objCtrl.InjectInportToStruct(templateWithData)
		if err != nil {
			return nil, err
		}

		// reformat outport.go
		err = r.outport.Reformat(ctx, objCtrl.GetControllerRouterFileName(), bytes)
		if err != nil {
			return nil, err
		}
	}

	// inject router for register
	//func (r *Controller) RegisterRouter() {
	//  r.Router.POST("/createorder", r.authorized(), r.createOrderHandler(r.CreateOrderInport)) <-- here
	//}
	{
		templateCode := r.outport.GetRouterRegisterTemplate(ctx, req.DriverName)

		templateWithData, err := r.outport.PrintTemplate(ctx, templateCode, objDataCtrl)
		if err != nil {
			return nil, err
		}

		bytes, err := objCtrl.InjectRouterBind(templateWithData)
		if err != nil {
			return nil, err
		}

		// reformat outport.go
		err = r.outport.Reformat(ctx, objCtrl.GetControllerRouterFileName(), bytes)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

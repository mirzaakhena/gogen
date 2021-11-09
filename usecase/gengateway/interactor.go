package gengateway

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/service"

	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/vo"
)

//go:generate mockery --name Outport -output mocks/

type genGatewayInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenGateway
func NewUsecase(outputPort Outport) Inport {
	return &genGatewayInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenGateway
func (r *genGatewayInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	packagePath := r.outport.GetPackagePath(ctx)

	err := service.CreateEverythingExactly("default/", "infrastructure/log", map[string]string{}, struct{PackagePath string}{PackagePath: packagePath})
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("default/", "infrastructure/util", map[string]string{}, struct{PackagePath string}{PackagePath: packagePath})
	if err != nil {
		return nil, err
	}

	obj, err := entity.NewObjGateway(req.GatewayName)
	if err != nil {
		return nil, err
	}

	notExistingMethod, err := r.createGatewayImpl(req.UsecaseName, packagePath, obj)
	if err != nil {
		return nil, err
	}

	gatewayCode := r.outport.GetGatewayMethodTemplate(ctx)

	// we will only inject the non existing method
	data := obj.GetData(packagePath, notExistingMethod)

	templateHasBeenInjected, err := r.outport.PrintTemplate(ctx, gatewayCode, data)
	if err != nil {
		return nil, err
	}

	bytes, err := obj.InjectToGateway(templateHasBeenInjected)
	if err != nil {
		return nil, err
	}

	// reformat outport.go
	err = r.outport.Reformat(ctx, obj.GetGatewayFileName(), bytes)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (r *genGatewayInteractor) createGatewayImpl(usecaseName string, packagePath string, obj *entity.ObjGateway) (vo.OutportMethods, error) {
	outportMethods, err := vo.NewOutportMethods(usecaseName, packagePath)
	if err != nil {
		return nil, err
	}

	err = service.CreateEverythingExactly("default/", "gateway", map[string]string{
		"gatewayname": obj.GatewayName.LowerCase(),
	}, obj.GetData(packagePath, outportMethods))
	if err != nil {
		return nil, err
	}

	// file gateway impl file is already exist, we want to inject non existing method
	existingFunc, err := vo.NewOutportMethodImpl("gateway", obj.GetGatewayRootFolderName(), packagePath)
	if err != nil {
		return nil, err
	}

	// collect the only methods that has not added yet
	notExistingMethod := vo.OutportMethods{}
	for _, m := range outportMethods {
		if _, exist := existingFunc[m.MethodName]; !exist {
			notExistingMethod = append(notExistingMethod, m)
		}
	}
	return notExistingMethod, nil
}

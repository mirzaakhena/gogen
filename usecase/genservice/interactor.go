package genservice

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
)

//go:generate mockery --name Outport -output mocks/

type genServiceInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenService
func NewUsecase(outputPort Outport) Inport {
	return &genServiceInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenService
func (r *genServiceInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// create object service
	obj, err := entity.NewObjService(req.ServiceName, req.UsecaseName)
	if err != nil {
		return nil, err
	}

	packagePath := r.outport.GetPackagePath(ctx)

	// create service
	err = service.CreateEverythingExactly("default/", "domain/service", map[string]string{}, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	// service._go file is already exist, but is the interface is exist ?
	exist, err := obj.IsServiceExist()
	if err != nil {
		return nil, err
	}

	if !exist {
		// check the prefix and give specific template for it
		templateCode := r.outport.GetServiceInterfaceTemplate(ctx)
		templateHasBeenInjected, err := r.outport.PrintTemplate(ctx, templateCode, obj.GetData(packagePath))
		if err != nil {
			return nil, err
		}

		bytes, err := obj.InjectCode(templateHasBeenInjected)
		if err != nil {
			return nil, err
		}

		// reformat interactor._go
		err = r.outport.Reformat(ctx, "domain/service/service._go", bytes)
		if err != nil {
			return nil, err
		}
	}

	// if usecase name is not empty means we want to inject to usecase
	if obj.ObjUsecase.UsecaseName.IsEmpty() {
		return res, nil
	}

	// inject to outport
	{
		// TODO declared twice in service
		err := obj.InjectToOutport()
		if err != nil {
			return nil, err
		}

		outportFile := obj.ObjUsecase.GetOutportFileName()

		// reformat outport._go
		err = r.outport.Reformat(ctx, outportFile, nil)
		if err != nil {
			return nil, err
		}
	}

	// inject to interactor
	{
		// check the prefix and give specific template for it
		interactorCode := r.outport.GetServiceInjectTemplate(ctx)
		templateHasBeenInjected, err := r.outport.PrintTemplate(ctx, interactorCode, obj.GetData(packagePath))
		if err != nil {
			return nil, err
		}

		interactorBytes, err := obj.InjectToInteractor(templateHasBeenInjected)
		if err != nil {
			return nil, err
		}

		// reformat interactor._go
		err = r.outport.Reformat(ctx, obj.ObjUsecase.GetInteractorFileName(), interactorBytes)
		if err != nil {
			return nil, err
		}
	}

	return res, nil
}

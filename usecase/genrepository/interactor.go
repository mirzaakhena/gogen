package genrepository

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
	"github.com/mirzaakhena/gogen/usecase/genentity"
)

//go:generate mockery --name Outport -output mocks/

type genRepositoryInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenRepository
func NewUsecase(outputPort Outport) Inport {
	return &genRepositoryInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenRepository
func (r *genRepositoryInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// create object repository
	obj, err := entity.NewObjRepository(req.RepositoryName, req.EntityName, req.UsecaseName)
	if err != nil {
		return nil, err
	}

	// create entity
	_, err = genentity.NewUsecase(r.outport).Execute(ctx, genentity.InportRequest{EntityName: req.EntityName})
	if err != nil {
		return nil, err
	}

	packagePath := r.outport.GetPackagePath(ctx)

	// create repository
	err = service.CreateEverythingExactly("default/", "domain/repository", map[string]string{}, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	// create service
	err = service.CreateEverythingExactly("default/", "domain/service", map[string]string{}, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	// repository._go file is already exist, but is the interface is exist ?
	exist, err := obj.IsRepoExist()
	if err != nil {
		return nil, err
	}

	if !exist {
		// check the prefix and give specific template for it
		templateCode := r.outport.GetRepoInterfaceTemplate(ctx, obj.RepositoryName)
		templateHasBeenInjected, err := r.outport.PrintTemplate(ctx, templateCode, obj.GetData(packagePath))
		if err != nil {
			return nil, err
		}

		bytes, err := obj.InjectCode(templateHasBeenInjected)
		if err != nil {
			return nil, err
		}

		// reformat interactor._go
		err = r.outport.Reformat(ctx, "domain/repository/repository._go", bytes)
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
		interactorCode := r.outport.GetRepoInjectTemplate(ctx, obj.RepositoryName)
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

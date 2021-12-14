package genusecase

import (
	"context"
	"github.com/mirzaakhena/gogen/domain/entity"
	"github.com/mirzaakhena/gogen/domain/service"
)

//go:generate mockery --name Outport -output mocks/

type genUsecaseInteractor struct {
	outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenUsecase
func NewUsecase(outputPort Outport) Inport {
	return &genUsecaseInteractor{
		outport: outputPort,
	}
}

// Execute the usecase GenUsecase
func (r *genUsecaseInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

	res := &InportResponse{}

	// create object usecase
	obj, err := entity.NewObjUsecase(req.UsecaseName)
	if err != nil {
		return nil, err
	}

	packagePath := r.outport.GetPackagePath(ctx)

	fileRenamer := map[string]string{
		"usecasename": obj.UsecaseName.LowerCase(),
	}

	//err = service.CreateEverythingExactly("default/", "domain/repository", fileRenamer, obj.GetData(packagePath))
	//if err != nil {
	//  return nil, err
	//}
	//
	//err = service.CreateEverythingExactly("default/", "domain/service", fileRenamer, obj.GetData(packagePath))
	//if err != nil {
	//  return nil, err
	//}

	err = service.CreateEverythingExactly("default/", "usecase", fileRenamer, obj.GetData(packagePath))
	if err != nil {
		return nil, err
	}

	//_, err = r.outport.CreateFolderIfNotExist(ctx, "domain/entity")
	//if err != nil {
	//  return nil, err
	//}
	//
	//_, err = r.outport.CreateFolderIfNotExist(ctx, "domain/vo")
	//if err != nil {
	//  return nil, err
	//}

	//// create file interactor._go
	//{
	//	tmp := r.outport.GetInteractorTemplate(ctx)
	//	outputFile := obj.GetInteractorFileName()
	//	_, err = r.outport.WriteFileIfNotExist(ctx, tmp, outputFile, obj.GetData(packagePath))
	//	if err != nil {
	//		return nil, err
	//	}
	//
	//  // reformat interactor._go since we have some import on it
	//  err = r.outport.Reformat(ctx, outputFile, nil)
	//  if err != nil {
	//    return nil, err
	//  }
	//}

	return res, nil
}

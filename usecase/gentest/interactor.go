package gentest

import (
  "context"
  "fmt"
  "github.com/mirzaakhena/gogen/domain/entity"
  "github.com/mirzaakhena/gogen/domain/service"
  "github.com/mirzaakhena/gogen/domain/vo"
)

//go:generate mockery --name Outport -output mocks/

type genTestInteractor struct {
  outport Outport
}

// NewUsecase is constructor for create default implementation of usecase GenTest
func NewUsecase(outputPort Outport) Inport {
  return &genTestInteractor{
    outport: outputPort,
  }
}

// Execute the usecase GenTest
func (r *genTestInteractor) Execute(ctx context.Context, req InportRequest) (*InportResponse, error) {

  res := &InportResponse{}

  packagePath := r.outport.GetPackagePath(ctx)

  err := service.CreateEverythingExactly("default/", "infrastructure/log", map[string]string{}, struct{PackagePath string}{PackagePath: packagePath})
  if err != nil {
    return nil, err
  }

  objUsecase, err := entity.NewObjUsecase(req.UsecaseName)
  if err != nil {
    return nil, err
  }

  obj, err := entity.NewObjTesting(req.TestName, *objUsecase)
  if err != nil {
    return nil, err
  }

  outportMethods, err := vo.NewOutportMethods(req.UsecaseName, packagePath)
  if err != nil {
    return nil, err
  }

  fileRenamer := map[string]string{
    "usecasename": objUsecase.UsecaseName.LowerCase(),
    "testname":    obj.TestName.LowerCase(),
  }
  err = service.CreateEverythingExactly("default/test/", "usecase", fileRenamer, obj.GetData(packagePath, outportMethods))
  if err != nil {
    return nil, err
  }

  // file gateway impl file is already exist, we want to inject non existing method
  existingFunc, err := vo.NewOutportMethodImpl(fmt.Sprintf("mockOutport%s", obj.TestName.PascalCase()), fmt.Sprintf("usecase/%s", obj.ObjUsecase.UsecaseName.LowerCase()), packagePath)
  if err != nil {
    return nil, err
  }

  fmt.Printf("%v\n", existingFunc)

  // collect the only methods that has not added yet
  notExistingMethod := vo.OutportMethods{}
  for _, m := range outportMethods {
    if _, exist := existingFunc[m.MethodName]; !exist {
      notExistingMethod = append(notExistingMethod, m)
    }
  }

  gatewayCode := r.outport.GetTestMethodTemplate(ctx)

  // we will only inject the non existing method
  data := obj.GetData(packagePath, notExistingMethod)

  templateHasBeenInjected, err := r.outport.PrintTemplate(ctx, gatewayCode, data)
  if err != nil {
    return nil, err
  }

  bytes, err := obj.InjectToTest(templateHasBeenInjected)
  if err != nil {
    return nil, err
  }

  // reformat outport.go
  err = r.outport.Reformat(ctx, obj.GetTestFileName(), bytes)
  if err != nil {
    return nil, err
  }

  return res, nil
}

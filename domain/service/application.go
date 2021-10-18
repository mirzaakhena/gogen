package service

import (
  "context"
  "fmt"
)

// ApplicationActionInterface ...
type ApplicationActionInterface interface {
  CreateFolderIfNotExistService
  WriteFileIfNotExistService
  GetErrorTemplate(ctx context.Context) (string, string)
  GetConstantTemplate(ctx context.Context) string
  GetApplicationTemplate(ctx context.Context) string
}

// ConstructApplication ...
func ConstructApplication(ctx context.Context, packagePath string, action ApplicationActionInterface) error {

  {
    //_, _ = action.CreateFolderIfNotExist(ctx, entity.GetErrorRootFolderName())
    _, _ = action.CreateFolderIfNotExist(ctx, "application/constant")
    _, _ = action.CreateFolderIfNotExist(ctx, "application/registry")
  }

  //errorFunc, errorEnum := action.GetErrorTemplate(ctx)
  //{
  //  outputFile := entity.GetErrorEnumFileName()
  //  _, err := action.WriteFileIfNotExist(ctx, errorEnum, outputFile, struct{}{})
  //  if err != nil {
  //    return err
  //  }
  //}
  //
  //{
  //  outputFile := entity.GetErrorFuncFileName()
  //  _, err := action.WriteFileIfNotExist(ctx, errorFunc, outputFile, struct{}{})
  //  if err != nil {
  //    return err
  //  }
  //}

  {
    appTemplateFile := action.GetConstantTemplate(ctx)
    outputFile := fmt.Sprintf("application/constant/constant.go")
    _, err := action.WriteFileIfNotExist(ctx, appTemplateFile, outputFile, struct{}{})
    if err != nil {
      return err
    }
  }

  {
    appTemplateFile := action.GetApplicationTemplate(ctx)
    outputFile := fmt.Sprintf("application/application.go")
    _, err := action.WriteFileIfNotExist(ctx, appTemplateFile, outputFile, struct {
      PackagePath string
    }{PackagePath: packagePath})
    if err != nil {
      return err
    }
  }

  return nil
}

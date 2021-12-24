package service

import (
	"context"
	"github.com/mirzaakhena/gogen/model/entity"
)

// ErrorActionInterface ...
type ErrorActionInterface interface {
	CreateFolderIfNotExistService
	WriteFileIfNotExistService
	GetErrorEnumTemplate(ctx context.Context) string
	GetErrorFuncTemplate(ctx context.Context) string
}

// ConstructError ...
func ConstructError(ctx context.Context, action ErrorActionInterface) error {

	{
		folder := entity.ObjError{}.GetErrorRootFolderName()
		_, err := action.CreateFolderIfNotExist(ctx, folder)
		if err != nil {
			return err
		}
	}

	{
		logTemplateFile := action.GetErrorEnumTemplate(ctx)
		outputFile := entity.ObjError{}.GetErrorEnumFileName()
		_, err := action.WriteFileIfNotExist(ctx, logTemplateFile, outputFile, struct{}{})
		if err != nil {
			return err
		}
	}

	{
		logImplTemplateFile := action.GetErrorFuncTemplate(ctx)
		outputFile := entity.ObjError{}.GetErrorFuncFileName()
		_, err := action.WriteFileIfNotExist(ctx, logImplTemplateFile, outputFile, struct{}{})
		if err != nil {
			return err
		}
	}

	return nil
}

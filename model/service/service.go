package service

import (
  "context"
)

// CreateFolderIfNotExistService ...
type CreateFolderIfNotExistService interface {
  CreateFolderIfNotExist(ctx context.Context, folderPath string) (bool, error)
}

// WriteFileIfNotExistService ...
type WriteFileIfNotExistService interface {
  WriteFileIfNotExist(ctx context.Context, templateFile, outputFilePath string, obj interface{}) (bool, error)
}

// WriteFileService ...
type WriteFileService interface {
  WriteFile(ctx context.Context, templateFile, outputFilePath string, obj interface{}) error
}

// ReformatService ...
type ReformatService interface {
  Reformat(ctx context.Context, goFilename string, bytes []byte) error
}

// GetPackagePathService ...
type GetPackagePathService interface {
  GetPackagePath(ctx context.Context) string
}

// IsFileExistService ...
type IsFileExistService interface {
  IsFileExist(ctx context.Context, filepath string) bool
}

// PrintTemplateService ...
type PrintTemplateService interface {
  PrintTemplate(ctx context.Context, templateString string, x interface{}) (string, error)
}
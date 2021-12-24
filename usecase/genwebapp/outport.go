package genwebapp

import (
	"context"
	"github.com/mirzaakhena/gogen/model/service"
)

// Outport of usecase
type Outport interface {
	service.GetPackagePathService
	service.WriteFileIfNotExistService

	GetWebappGitignoreTemplate(ctx context.Context) string
	GetWebappIndexTemplate(ctx context.Context) string
	GetWebappPackageTemplate(ctx context.Context) string
	GetWebappReadmeTemplate(ctx context.Context) string
	GetWebappViteConfigTemplate(ctx context.Context) string
}

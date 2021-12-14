package prod

import (
	"context"
	"github.com/mirzaakhena/gogen/infrastructure/templates"
)

func (r *prodGateway) GetWebappGitignoreTemplate(ctx context.Context) string {
	return templates.ReadFileCrud("._gitignore")
}

func (r *prodGateway) GetWebappIndexTemplate(ctx context.Context) string {
	return templates.ReadFileCrud("index._html")
}

func (r *prodGateway) GetWebappPackageTemplate(ctx context.Context) string {
	return templates.ReadFileCrud("package._json")
}

func (r *prodGateway) GetWebappReadmeTemplate(ctx context.Context) string {
	return templates.ReadFileCrud("README._md")
}

func (r *prodGateway) GetWebappViteConfigTemplate(ctx context.Context) string {
	return templates.ReadFileCrud("vite.config._js")
}

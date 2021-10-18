package prod

import (
  "context"
  "github.com/mirzaakhena/gogen/infrastructure/templates"
)

//// GetErrorTemplate ...
//func (r prodGateway) GetErrorTemplate(ctx context.Context) (fun string, enm string) {
//	return templates.ApplicationErrorFuncFile, templates.ApplicationErrorEnumFile
//}
//
//// GetConstantTemplate ...
//func (r prodGateway) GetConstantTemplate(ctx context.Context) string {
//	return templates.ApplicationConstantTemplateFile
//}
//
//// GetApplicationTemplate ...
//func (r prodGateway) GetApplicationTemplate(ctx context.Context) string {
//	return templates.ApplicationFile
//}
//
// GetErrorLineTemplate ...
func (r prodGateway) GetErrorLineTemplate(ctx context.Context) string {
	return templates.ReadFile("application/apperror/~inject._go")
}

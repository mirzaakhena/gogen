package templates

import (
	"embed"
	_ "embed"
	"fmt"
) // used for embed file

//go:embed default
var AppTemplates embed.FS

//go:embed crud
var CrudTemplates embed.FS

//go:embed webapp
var WebappTemplates embed.FS

func ReadFile(filepath string) string {
	bytes, err := AppTemplates.ReadFile(fmt.Sprintf("default/%s", filepath))
	if err != nil {
		panic(err.Error())
	}
	return string(bytes)
}

func ReadFileCrud(filepath string) string {
	bytes, err := CrudTemplates.ReadFile(fmt.Sprintf("crud/%s", filepath))
	if err != nil {
		panic(err.Error())
	}
	return string(bytes)
}

func ReadFileWebapp(filepath string) string {
	bytes, err := WebappTemplates.ReadFile(fmt.Sprintf("webapp/%s", filepath))
	if err != nil {
		panic(err.Error())
	}
	return string(bytes)
}

var (

	// LogFile ...
	//go:embed default/infrastructure/log/log._go
	LogFile string

	// LogDefaultFile ...
	//go:embed default/infrastructure/log/log_default._go
	LogDefaultFile string

	// EntityFile ...
	//go:embed default/domain/entity/${entityname}._go
	EntityFile string
)

var (

	//go:embed default/application/registry/~driver_gin._go
	RegistryGingonicFile string
)

var (

	//go:embed default/_controller/gingonic/handler-func._go
	ControllerGinGonicHandlerFuncFile string

	//go:embed default/_controller/gingonic/interceptor._go
	ControllerGinGonicInterceptorFile string

	//go:embed default/_controller/gingonic/response._go
	ControllerGinGonicResponseFile string

	//go:embed default/_controller/gingonic/router-inport._go
	ControllerGinGonicRouterInportFile string

	//go:embed default/_controller/gingonic/router-struct._go
	ControllerGinGonicRouterStructFile string

	//go:embed default/_controller/gingonic/router-register._go
	ControllerGinGonicRouterRegisterFile string
)

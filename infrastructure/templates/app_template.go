package templates

import (
	"embed"
	_ "embed"
	"fmt"
) // used for embed file

//go:embed default
var AppTemplates embed.FS

func ReadFile(filepath string) string {
	bytes, err := AppTemplates.ReadFile(fmt.Sprintf("default/%s", filepath))
	if err != nil {
		panic(err.Error())
	}
	return string(bytes)
}

var (

	//// InportFile ...
	////go:embed default/_usecase/usecase_inport._go
	//InportFile string
	//
	//// OutportFile ...
	////go:embed default/_usecase/usecase_outport._go
	//OutportFile string
	//
	//// InteractorFile ...
	////go:embed default/_usecase/usecase_interactor._go
	//InteractorFile string
	//
	//// TestFile ...
	////go:embed default/_usecase/usecase_test._go
	//TestFile string

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

//// RepositoryFile ...
////go:embed default/domain/repository/repository._go
//RepositoryFile string
//
//// RepositoryInterfaceFile ...
////go:embed default/domain/repository/repository_interface._go
//RepositoryInterfaceFile string
//
//// RepositoryInterfaceFindFile ...
////go:embed default/domain/repository/repository_interface_find._go
//RepositoryInterfaceFindFile string
//
//// RepositoryInterfaceFindOneFile ...
////go:embed default/domain/repository/repository_interface_findone._go
//RepositoryInterfaceFindOneFile string
//
//// RepositoryInterfaceRemoveFile ...
////go:embed default/domain/repository/repository_interface_remove._go
//RepositoryInterfaceRemoveFile string
//
//// RepositoryInterfaceSaveFile ...
////go:embed default/domain/repository/repository_interface_save._go
//RepositoryInterfaceSaveFile string
)

var (

//// RepoInjectInteractorFile ...
////go:embed default/domain/repository/repository_inject._go
//RepoInjectInteractorFile string
//
//// RepoInjectInteractorFindFile ...
////go:embed default/domain/repository/repository_inject_find._go
//RepoInjectInteractorFindFile string
//
//// RepoInjectInteractorFindOneFile ...
////go:embed default/domain/repository/repository_inject_findone._go
//RepoInjectInteractorFindOneFile string
//
//// RepoInjectInteractorSaveFile ...
////go:embed default/domain/repository/repository_inject_save._go
//RepoInjectInteractorSaveFile string
//
//// RepoInjectInteractorRemoveFile ...
////go:embed default/domain/repository/repository_inject_remove._go
//RepoInjectInteractorRemoveFile string
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

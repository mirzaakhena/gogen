package templates

import _ "embed"

//go:embed inport._go
var InportFile string

//go:embed outport._go
var OutportFile string

//go:embed interactor._go
var InteractorFile string

//go:embed entity._go
var EntityFile string

//go:embed method._go
var MethodFile string

//go:embed enum._go
var EnumFile string

//go:embed error_enum._go
var ErrorEnumFile string

//go:embed error_func._go
var ErrorFuncFile string

//go:embed error_enum_template._go
var ErrorEnumTemplateFile string

//go:embed valueobject._go
var ValueObjectFile string

//go:embed valuestring._go
var ValueStringFile string

//go:embed repository._go
var RepositoryFile string

//go:embed repository_interface._go
var RepositoryInterfaceFile string

//go:embed repository_interface_remove._go
var RepositoryRemoveFile string

//go:embed repository_interface_save._go
var RepositorySaveFile string

//go:embed repository_interface_find._go
var RepositoryFindFile string

//go:embed repository_interface_findone._go
var RepositoryFindOneFile string

//go:embed repository_inject_find._go
var RepositoryInjectFindFile string

//go:embed repository_inject_findone._go
var RepositoryInjectFindOneFile string

//go:embed repository_inject_remove._go
var RepositoryInjectRemoveFile string

//go:embed repository_inject_save._go
var RepositoryInjectSaveFile string

//go:embed repository_inject._go
var RepositoryInjectInterfaceFile string

//go:embed service._go
var ServiceFile string

//go:embed service_inject._go
var ServiceInjectFile string

//go:embed service_template._go
var ServiceTemplateFile string

//go:embed gateway._go
var GatewayFile string

//go:embed controller._go
var ControllerFile string

//go:embed controller-gin._go
var ControllerGinFile string

//go:embed controller_response._go
var ControllerResponseFile string

//go:embed infra_log_contract._go
var LogContractFile string

//go:embed infra_log_impl._go
var LogImplFile string

//go:embed infra_log_public._go
var LogPublicFile string

//go:embed infra_helper._go
var HelperFile string

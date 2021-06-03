package templates

import _ "embed"

//go:embed usecase/usecase_inport._go
var InportFile string

//go:embed usecase/usecase_outport._go
var OutportFile string

//go:embed usecase/usecase_interactor._go
var InteractorFile string

// Entity
var (
	//go:embed entity/domain_entity._go
	EntityFile string

	//go:embed entity/domain_entity_method._go
	MethodFile string

	//go:embed entity/domain_vo_enum._go
	EnumFile string
)

// VO
var (
	//go:embed vo/domain_vo_valueobject._go
	ValueObjectFile string

	//go:embed vo/domain_vo_valuestring._go
	ValueStringFile string
)

//go:embed error/error_enum._go
var ErrorEnumFile string

//go:embed error/error_func._go
var ErrorFuncFile string

//go:embed error/error_enum_template._go
var ErrorEnumTemplateFile string

// Repository
var (

	//go:embed repository/repository._go
	RepositoryFile string

	//go:embed repository/repository_interface._go
	RepositoryInterfaceFile string

	//go:embed repository/repository_interface_remove._go
	RepositoryRemoveFile string

	//go:embed repository/repository_interface_save._go
	RepositorySaveFile string

	//go:embed repository/repository_interface_find._go
	RepositoryFindFile string

	//go:embed repository/repository_interface_findone._go
	RepositoryFindOneFile string

	//go:embed repository/repository_inject_find._go
	RepositoryInjectFindFile string

	//go:embed repository/repository_inject_findone._go
	RepositoryInjectFindOneFile string

	//go:embed repository/repository_inject_remove._go
	RepositoryInjectRemoveFile string

	//go:embed repository/repository_inject_save._go
	RepositoryInjectSaveFile string

	//go:embed repository/repository_inject._go
	RepositoryInjectInterfaceFile string
)

//go:embed service/service._go
var ServiceFile string

//go:embed service/service_inject._go
var ServiceInjectFile string

//go:embed service/service_template._go
var ServiceTemplateFile string

//go:embed log/infra_log._go
var LogContractFile string

//go:embed log/infra_log_default._go
var LogImplFile string

//go:embed infra_helper._go
var HelperFile string

//go:embed infra_server_shutdown._go
var ServerShutdownFile string

//go:embed application._go
var ApplicationFile string

//go:embed main._go
var MainFile string

//go:embed usecase/usecase_test._go
var TestFile string

// Gateway
var (
	//go:embed gateway/gateway._go
	GatewayFile string

	//go:embed gateway/gateway_inject_method._go
	GatewayMethodFile string
)

//go:embed README._md
var ReadmeFile string
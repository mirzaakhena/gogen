package gogen

type applicationSchema struct {
}

func NewApplicationSchema() Generator {
	return &applicationSchema{}
}

func (d *applicationSchema) Generate(args ...string) error {

	CreateFolder(".application_schema/usecases")

	// CreateFolder(".application_schema/models")

	// CreateFolder(".application_schema/controllers")

	// CreateFolder(".application_schema/datasources")

	CreateFolder("binder/")

	CreateFolder("controllers/")

	CreateFolder("datasources/mocks")

	CreateFolder("entities/model")

	CreateFolder("entities/repository")

	CreateFolder("shared/")

	CreateFolder("binder/")

	CreateFolder("inport/")

	CreateFolder("interactor/")

	CreateFolder("outport/")

	CreateFolder("utils/")

	mainApps := MainApps{PackagePath: GetPackagePath()}

	WriteFileIfNotExist(
		"main._go",
		"main.go",
		&mainApps,
	)

	WriteFileIfNotExist(
		"config._toml",
		"config.toml",
		struct{}{},
	)

	WriteFileIfNotExist(
		"README._md",
		"README.md",
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/runner._go",
		"binder/runner.go",
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/setup_component._go",
		"binder/setup_component.go",
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/wiring_component._go",
		"binder/wiring_component.go",
		struct{}{},
	)

	return nil
}

type MainApps struct {
	PackagePath string ``
}

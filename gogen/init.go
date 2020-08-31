package gogen

type applicationSchema struct {
}

func NewApplicationSchema() Generator {
	return &applicationSchema{}
}

func (d *applicationSchema) Generate(args ...string) error {

	CreateFolder(".application_schema/usecases")

	CreateFolder(".application_schema/models")

	CreateFolder(".application_schema/controllers")

	CreateFolder(".application_schema/datasources")

	CreateFolder("binder/")

	CreateFolder("controllers/")

	CreateFolder("datasources/mocks")

	CreateFolder("entities/model")

	CreateFolder("entities/repository")

	CreateFolder("shared/")

	CreateFolder("binder/")

	CreateFolder("usecases/")

	CreateFolder("utils/")

	WriteFile(
		"main._go",
		"main.go",
		struct{}{},
	)

	WriteFile(
		"config._toml",
		"config.toml",
		struct{}{},
	)

	WriteFile(
		"README._md",
		"README.md",
		struct{}{},
	)

	return nil
}

package gogen

import (
	"fmt"
	"strings"
)

type applicationSchema struct {
}

func NewApplicationSchema() Generator {
	return &applicationSchema{}
}

func (d *applicationSchema) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("try `gogen init .` for current directory as active project\nor  `gogen init <your project name>` for create new project directory. But remember do `cd <your project name>` to start working")
	}

	baseFolder := fmt.Sprintf("./%s/", strings.TrimSpace(args[2]))

	CreateFolder("%s.application_schema/usecases", baseFolder)

	CreateFolder("%sbinder/", baseFolder)

	CreateFolder("%scontrollers/", baseFolder)

	CreateFolder("%sdatasources/mocks", baseFolder)

	CreateFolder("%sentities/model", baseFolder)

	CreateFolder("%sentities/repository", baseFolder)

	CreateFolder("%sshared/", baseFolder)

	CreateFolder("%sbinder/", baseFolder)

	CreateFolder("%sinport/", baseFolder)

	CreateFolder("%sinteractor/", baseFolder)

	CreateFolder("%soutport/", baseFolder)

	CreateFolder("%sutils/", baseFolder)

	WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%smain.go", baseFolder),
		struct {
			PackagePath string
		}{
			PackagePath: GetPackagePath(),
		},
	)

	WriteFileIfNotExist(
		"config._toml",
		fmt.Sprintf("%sconfig.toml", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"README._md",
		fmt.Sprintf("%sREADME.md", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/runner._go",
		fmt.Sprintf("%sbinder/runner.go", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/setup_component._go",
		fmt.Sprintf("%sbinder/setup_component.go", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"binder/wiring_component._go",
		fmt.Sprintf("%sbinder/wiring_component.go", baseFolder),
		struct{}{},
	)

	return nil
}

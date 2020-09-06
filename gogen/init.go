package gogen

import (
	"fmt"
	"strings"
)

const (
	INIT_PROJECT_INDEX int = 2
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

	directory := strings.TrimSpace(args[2])

	baseFolder := fmt.Sprintf("%s/", directory)

	folder := ""
	if directory != "." {
		folder = fmt.Sprintf("/%s", strings.TrimSpace(args[2]))
	}

	fmt.Printf("directory :%s\n", folder)

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
			Directory   string
		}{
			PackagePath: GetPackagePath(),
			Directory:   folder,
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
		struct {
			PackagePath string
		}{
			PackagePath: GetPackagePath(),
		},
	)

	return nil
}

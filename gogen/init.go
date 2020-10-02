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

	directory := strings.TrimSpace(args[2])

	baseFolder := fmt.Sprintf("%s/", directory)

	folder := ""
	if directory != "." {
		folder = fmt.Sprintf("/%s", strings.TrimSpace(args[2]))
	}

	CreateFolder("%sapplication/", baseFolder)

	CreateFolder("%scontroller/", baseFolder)

	CreateFolder("%sdatasource/", baseFolder)

	CreateFolder("%smodel/", baseFolder)

	CreateFolder("%susecase/", baseFolder)

	CreateFolder("%sutil/", baseFolder)

	_ = WriteFileIfNotExist(
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

	_ = WriteFileIfNotExist(
		"config._toml",
		fmt.Sprintf("%sconfig.toml", baseFolder),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"README._md",
		fmt.Sprintf("%sREADME.md", baseFolder),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/runner._go",
		fmt.Sprintf("%sapplication/runner.go", baseFolder),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/application._go",
		fmt.Sprintf("%sapplication/application.go", baseFolder),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/schema._go",
		fmt.Sprintf("%sapplication/schema.go", baseFolder),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/registry._go",
		fmt.Sprintf("%sapplication/registry.go", baseFolder),
		struct{}{},
	)

	return nil
}

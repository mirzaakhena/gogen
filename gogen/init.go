package gogen

import (
	"fmt"
	"strings"

	"github.com/mirzaakhena/templator"
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

	templator.CreateFolder("%sapplication/", baseFolder)

	templator.CreateFolder("%scontroller/", baseFolder)

	templator.CreateFolder("%sdatasource/", baseFolder)

	templator.CreateFolder("%smodel/", baseFolder)

	templator.CreateFolder("%susecase/", baseFolder)

	templator.CreateFolder("%sutil/", baseFolder)

	_ = templator.WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%smain.go", baseFolder),
		struct {
			PackagePath string
			Directory   string
		}{
			PackagePath: templator.GetPackagePath(),
			Directory:   folder,
		},
	)

	_ = templator.WriteFileIfNotExist(
		"config._toml",
		fmt.Sprintf("%sconfig.toml", baseFolder),
		struct{}{},
	)

	_ = templator.WriteFileIfNotExist(
		"README._md",
		fmt.Sprintf("%sREADME.md", baseFolder),
		struct{}{},
	)

	_ = templator.WriteFileIfNotExist(
		"application/runner._go",
		fmt.Sprintf("%sapplication/runner.go", baseFolder),
		struct{}{},
	)

	_ = templator.WriteFileIfNotExist(
		"application/application._go",
		fmt.Sprintf("%sapplication/application.go", baseFolder),
		struct{}{},
	)

	_ = templator.WriteFileIfNotExist(
		"application/schema._go",
		fmt.Sprintf("%sapplication/schema.go", baseFolder),
		struct{}{},
	)

	_ = templator.WriteFileIfNotExist(
		"application/registry._go",
		fmt.Sprintf("%sapplication/registry.go", baseFolder),
		struct{}{},
	)

	return nil
}

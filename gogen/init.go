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

	CreateFolder("%scontrollers/", baseFolder)

	CreateFolder("%sdatasources/", baseFolder)

	CreateFolder("%sentities/", baseFolder)

	CreateFolder("%susecases/", baseFolder)

	CreateFolder("%srepositories/", baseFolder)

	CreateFolder("%sservices/", baseFolder)

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
		"application/runner._go",
		fmt.Sprintf("%sapplication/runner.go", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"application/setup._go",
		fmt.Sprintf("%sapplication/setup.go", baseFolder),
		struct{}{},
	)

	WriteFileIfNotExist(
		"application/wiring._go",
		fmt.Sprintf("%sapplication/wiring.go", baseFolder),
		struct {
			PackagePath string
		}{
			PackagePath: GetPackagePath(),
		},
	)

	return nil
}

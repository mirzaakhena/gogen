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

	folderPath := strings.TrimSpace(args[2])

	return GenerateInit(folderPath)
}

func GenerateInit(folderPath string) error {

	var folderImport string
	if folderPath != "." {
		folderImport = fmt.Sprintf("/%s", folderPath)
	}

	CreateFolder("%s/application/", folderPath)

	CreateFolder("%s/controller/", folderPath)

	CreateFolder("%s/datasource/", folderPath)

	CreateFolder("%s/model/", folderPath)

	CreateFolder("%s/usecase/", folderPath)

	CreateFolder("%s/util/", folderPath)

	_ = WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%s/main.go", folderPath),
		struct {
			PackagePath string
			Directory   string
		}{
			PackagePath: GetPackagePath(),
			Directory:   folderImport,
		},
	)

	_ = WriteFileIfNotExist(
		"config._toml",
		fmt.Sprintf("%s/config.toml", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"README._md",
		fmt.Sprintf("%s/README.md", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/runner._go",
		fmt.Sprintf("%s/application/runner.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/application._go",
		fmt.Sprintf("%s/application/application.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/schema._go",
		fmt.Sprintf("%s/application/schema.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/registry._go",
		fmt.Sprintf("%s/application/registry.go", folderPath),
		struct{}{},
	)

	return nil
}

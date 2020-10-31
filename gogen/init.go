package gogen

import (
	"fmt"
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

	return GenerateInit(InitRequest{
		FolderPath: args[2],
	})
}

type InitRequest struct {
	FolderPath string
}

func GenerateInit(req InitRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	CreateFolder("%s/application/", req.FolderPath)

	CreateFolder("%s/controller/", req.FolderPath)

	CreateFolder("%s/datasource/", req.FolderPath)

	CreateFolder("%s/model/", req.FolderPath)

	CreateFolder("%s/usecase/", req.FolderPath)

	CreateFolder("%s/util/", req.FolderPath)

	_ = WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%s/main.go", req.FolderPath),
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
		fmt.Sprintf("%s/config.toml", req.FolderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"README._md",
		fmt.Sprintf("%s/README.md", req.FolderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/runner._go",
		fmt.Sprintf("%s/application/runner.go", req.FolderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/application._go",
		fmt.Sprintf("%s/application/application.go", req.FolderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/schema._go",
		fmt.Sprintf("%s/application/schema.go", req.FolderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/registry._go",
		fmt.Sprintf("%s/application/registry.go", req.FolderPath),
		struct{}{},
	)

	return nil
}

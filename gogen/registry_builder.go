package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type RegistryBuilderRequest struct {
	RegistryName   string
	UsecaseName    string
	DatasourceName string
	ControllerName string
	FolderPath     string
	Framework      string
}

type registryBuilder struct {
	RegistryBuilderRequest RegistryBuilderRequest
}

func NewRegistry(req RegistryBuilderRequest) Generator {
	return &registryBuilder{RegistryBuilderRequest: req}
}

func (d *registryBuilder) Generate() error {

	registryName := strings.TrimSpace(d.RegistryBuilderRequest.RegistryName)
	usecaseName := strings.TrimSpace(d.RegistryBuilderRequest.UsecaseName)
	datasourceName := strings.TrimSpace(d.RegistryBuilderRequest.DatasourceName)
	controllerName := strings.TrimSpace(d.RegistryBuilderRequest.ControllerName)
	folderPath := d.RegistryBuilderRequest.FolderPath
	framework := d.RegistryBuilderRequest.Framework

	if len(registryName) == 0 {
		return fmt.Errorf("Registry name must not empty")
	}

	if len(usecaseName) == 0 {
		return fmt.Errorf("Usecase name must not empty")
	}

	if len(datasourceName) == 0 {
		return fmt.Errorf("Datasource name must not empty")
	}

	if len(controllerName) == 0 {
		return fmt.Errorf("Controller name must not empty")
	}

	if !IsExist(fmt.Sprintf("%s/controller/%s/%s.go", folderPath, strings.ToLower(controllerName), usecaseName)) {
		return fmt.Errorf("controller %s/%s is not found", controllerName, usecaseName)
	}

	if !IsExist(fmt.Sprintf("%s/datasource/%s/%s.go", folderPath, strings.ToLower(datasourceName), usecaseName)) {
		return fmt.Errorf("datasource %s/%s is not found", datasourceName, usecaseName)
	}

	if !IsExist(fmt.Sprintf("%s/usecase/%s", folderPath, strings.ToLower(usecaseName))) {
		return fmt.Errorf("usecase %s is not found", usecaseName)
	}

	// create a folder with usecase name
	CreateFolder("%s/application/registry", folderPath)

	rg := StructureRegistry{
		RegistryName: registryName,
		PackagePath:  GetPackagePath(),
	}

	_ = WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%s/main.go", folderPath),
		rg,
	)

	_ = WriteFileIfNotExist(
		"application/application._go",
		fmt.Sprintf("%s/application/application.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/gracefully_shutdown._go",
		fmt.Sprintf("%s/application/gracefully_shutdown.go", folderPath),
		struct{}{},
	)

	var funcDeclareInjectedCode string

	funcCallInjectedCode, _ := PrintTemplate("application/registry/func_call._go", d.RegistryBuilderRequest)

	if framework == "nethttp" {
		_ = WriteFileIfNotExist(
			"application/handler_http._go",
			fmt.Sprintf("%s/application/handler_http.go", folderPath),
			struct{}{},
		)

		_ = WriteFileIfNotExist(
			"application/registry/registry_http._go",
			fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)),
			rg,
		)

		funcDeclareInjectedCode, _ = PrintTemplate("application/registry/func_http_declaration._go", d.RegistryBuilderRequest)

	} else //

	if framework == "gin" {
		_ = WriteFileIfNotExist(
			"application/handler_gin._go",
			fmt.Sprintf("%s/application/handler_gin.go", folderPath),
			struct{}{},
		)

		_ = WriteFileIfNotExist(
			"application/registry/registry_gin._go",
			fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)),
			rg,
		)

		funcDeclareInjectedCode, _ = PrintTemplate("application/registry/func_gin_declaration._go", d.RegistryBuilderRequest)

	} else //

	{
		return fmt.Errorf("not recognize framework")
	}

	// open registry file
	file, err := os.Open(fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)))
	if err != nil {
		return fmt.Errorf("not found registry file. You need to call 'gogen init .' first")
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		if strings.HasPrefix(strings.TrimSpace(row), "//code_injection function declaration") {
			buffer.WriteString(funcDeclareInjectedCode)
			buffer.WriteString("\n")
		} else //

		if strings.HasPrefix(strings.TrimSpace(row), "//code_injection function call") {
			buffer.WriteString(funcCallInjectedCode)
			buffer.WriteString("\n")
		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)), buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

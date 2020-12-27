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
}

type registryBuilder struct {
	RegistryBuilderRequest RegistryBuilderRequest
}

func NewRegistry(req RegistryBuilderRequest) Generator {
	return &registryBuilder{RegistryBuilderRequest: req}
}

func (d *registryBuilder) Generate() error {

	registryName := d.RegistryBuilderRequest.RegistryName
	usecaseName := d.RegistryBuilderRequest.UsecaseName
	datasourceName := d.RegistryBuilderRequest.DatasourceName
	controllerName := d.RegistryBuilderRequest.ControllerName
	folderPath := d.RegistryBuilderRequest.FolderPath

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
		"application/application._go",
		fmt.Sprintf("%s/application/application.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/gracefully_shutdown._go",
		fmt.Sprintf("%s/application/gracefully_shutdown.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/http_handler._go",
		fmt.Sprintf("%s/application/http_handler.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/registry/registry._go",
		fmt.Sprintf("%s/application/registry/%s.go", folderPath, registryName),
		rg,
	)

	_ = WriteFileIfNotExist(
		"application/registry/registry._go",
		fmt.Sprintf("%s/application/registry/%s.go", folderPath, registryName),
		rg,
	)

	_ = WriteFileIfNotExist(
		"main._go",
		fmt.Sprintf("%s/main.go", folderPath),
		rg,
	)

	// open registry file
	file, err := os.Open(fmt.Sprintf("%s/application/registry/%s.go", folderPath, registryName))
	if err != nil {
		return fmt.Errorf("not found registry file. You need to call 'gogen init .' first")
	}
	defer file.Close()

	funcCallInjectedCode, _ := PrintTemplate("application/registry/func_call._go", d.RegistryBuilderRequest)
	funcDeclareInjectedCode, _ := PrintTemplate("application/registry/func_declaration._go", d.RegistryBuilderRequest)

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

	if err := ioutil.WriteFile(fmt.Sprintf("%s/application/registry/%s.go", folderPath, registryName), buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

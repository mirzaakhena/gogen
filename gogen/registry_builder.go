package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"os"
	"strings"
)

type RegistryBuilderRequest struct {
	RegistryName   string
	UsecaseName    string
	GatewayName    string
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
	gatewayName := strings.TrimSpace(d.RegistryBuilderRequest.GatewayName)
	controllerName := strings.TrimSpace(d.RegistryBuilderRequest.ControllerName)
	folderPath := d.RegistryBuilderRequest.FolderPath
	framework := d.RegistryBuilderRequest.Framework

	packagePath := GetPackagePath()

	if len(registryName) == 0 {
		return fmt.Errorf("Registry name must not empty")
	}

	if len(usecaseName) == 0 {
		return fmt.Errorf("Usecase name must not empty")
	}

	if len(gatewayName) == 0 {
		return fmt.Errorf("Gateway name must not empty")
	}

	if len(controllerName) == 0 {
		return fmt.Errorf("Controller name must not empty")
	}

	if !IsExist(fmt.Sprintf("%s/controller/%s/%s.go", folderPath, strings.ToLower(controllerName), PascalCase(usecaseName))) {
		return fmt.Errorf("controller %s/%s is not found", controllerName, PascalCase(usecaseName))
	}

	if !IsExist(fmt.Sprintf("%s/gateway/%s/%s.go", folderPath, strings.ToLower(gatewayName), PascalCase(usecaseName))) {
		return fmt.Errorf("gateway %s/%s is not found", gatewayName, PascalCase(usecaseName))
	}

	if !IsExist(fmt.Sprintf("%s/usecase/%s", folderPath, strings.ToLower(usecaseName))) {
		return fmt.Errorf("usecase %s is not found", PascalCase(usecaseName))
	}

	// create a folder with usecase name
	CreateFolder("%s/application/registry", folderPath)

	rg := StructureRegistry{
		RegistryName: registryName,
		PackagePath:  packagePath,
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
			fmt.Sprintf("%s/application/registry/%s.go", folderPath, registryName),
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
			fmt.Sprintf("%s/application/registry/%s.go", folderPath, registryName),
			rg,
		)

		funcDeclareInjectedCode, _ = PrintTemplate("application/registry/func_gin_declaration._go", d.RegistryBuilderRequest)

	} else //

	{
		return fmt.Errorf("not recognize framework")
	}

	// open registry file

	registryFile := fmt.Sprintf("%s/application/registry/%s.go", folderPath, registryName)
	file, err := os.Open(registryFile)
	if err != nil {
		return fmt.Errorf("not found registry file. You need to call 'gogen init .' first")
	}
	defer file.Close()

	fSet := token.NewFileSet()
	node, errParse := parser.ParseFile(fSet, registryFile, nil, parser.ParseComments)
	if errParse != nil {
		return errParse
	}

	existingImportMap := ReadImports(node)

	scanner := bufio.NewScanner(file)

	importMode := false
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		if importMode && strings.HasPrefix(row, ")") {
			importMode = false

			if _, exist := existingImportMap[fmt.Sprintf("\"%s/controller/%s\"", packagePath, controllerName)]; !exist {
				buffer.WriteString(fmt.Sprintf("	\"%s/controller/%s\"", packagePath, controllerName))
				buffer.WriteString("\n")
			}

			if _, exist := existingImportMap[fmt.Sprintf("\"%s/gateway/%s\"", packagePath, gatewayName)]; !exist {
				buffer.WriteString(fmt.Sprintf("	\"%s/gateway/%s\"", packagePath, gatewayName))
				buffer.WriteString("\n")
			}

			if _, exist := existingImportMap[fmt.Sprintf("\"%s/usecase/%s\"", packagePath, strings.ToLower(usecaseName))]; !exist {
				buffer.WriteString(fmt.Sprintf("	\"%s/usecase/%s\"", packagePath, strings.ToLower(usecaseName)))
				buffer.WriteString("\n")
			}

		} else //

		if strings.HasPrefix(row, "import (") {
			importMode = true

		} else //

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

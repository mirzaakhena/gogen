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
	FolderPath     string
	GomodPath      string
	RegistryName   string
	UsecaseName    string
	GatewayName    string
	ControllerName string
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
	gomodPath := d.RegistryBuilderRequest.GomodPath

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

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

	CreateFolder("%s/infrastructure/config", folderPath)
	CreateFolder("%s/infrastructure/database", folderPath)
	CreateFolder("%s/infrastructure/server", folderPath)
	CreateFolder("%s/infrastructure/log", folderPath)

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
		"config-registryName._toml",
		fmt.Sprintf("%s/config-%s.toml", folderPath, LowerCase(registryName)),
		rg,
	)

	_ = WriteFileIfNotExist(
		"application/application._go",
		fmt.Sprintf("%s/application/application.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"infrastructure/config/config._go",
		fmt.Sprintf("%s/infrastructure/config/config.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"infrastructure/server/gracefully_shutdown._go",
		fmt.Sprintf("%s/infrastructure/server/gracefully_shutdown.go", folderPath),
		rg,
	)

	_ = WriteFileIfNotExist(
		"infrastructure/server/http_handler._go",
		fmt.Sprintf("%s/infrastructure/server/http_handler.go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"application/registry/registry._go",
		fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)),
		rg,
	)

	funcCallInjectedCode, _ := PrintTemplate("application/registry/func_call._go", d.RegistryBuilderRequest)
	funcDeclareInjectedCode, _ := PrintTemplate("application/registry/func_declaration._go", d.RegistryBuilderRequest)

	// open registry file

	registryFile := fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName))
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

	methodCallMode := false
	importMode := false
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		if methodCallMode {

			if strings.HasPrefix(row, "}") {
				methodCallMode = false
				buffer.WriteString(funcCallInjectedCode)
				buffer.WriteString("\n")
			}
			//else {
			// TODO detecting if any method has been called
			// }

		}

		if strings.HasPrefix(row, fmt.Sprintf("func (r *%sRegistry) RegisterUsecase() {", CamelCase(registryName))) {
			methodCallMode = true

		} else //

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

			if _, exist := existingImportMap[fmt.Sprintf("\"%s/usecase/%s\"", packagePath, LowerCase(usecaseName))]; !exist {
				buffer.WriteString(fmt.Sprintf("	\"%s/usecase/%s\"", packagePath, LowerCase(usecaseName)))
				buffer.WriteString("\n")
			}

		} else //

		if strings.HasPrefix(row, "import (") {
			importMode = true

		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	buffer.WriteString(funcDeclareInjectedCode)
	buffer.WriteString("\n")

	if err := ioutil.WriteFile(fmt.Sprintf("%s/application/registry/%s.go", folderPath, PascalCase(registryName)), buffer.Bytes(), 0644); err != nil {
		return err
	}

	GoFormat(packagePath)

	return nil
}

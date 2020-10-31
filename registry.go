package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type registry struct {
}

func NewRegistry() Generator {
	return &registry{}
}

const funcDeclare = `
func %sHandler(a *Application) {
	outport := %s.New%sDatasource()
	inport := %s.New%sUsecase(outport)
	a.Router.POST("/%s", %s.%s(inport))
}`

const funcCall = "	%sHandler(a)"

func (d *registry) Generate(args ...string) error {

	if len(args) < 5 {
		return fmt.Errorf("please define controller, datasource, usecase name. ex: `gogen registry restapi production CreateMenu`")
	}

	return GenerateRegistry(RegistryRequest{
		ControllerName: args[2],
		DatasourceName: args[3],
		UsecaseName:    args[4],
		FolderPath:     ".",
	})

}

type RegistryRequest struct {
	ControllerName string
	DatasourceName string
	UsecaseName    string
	FolderPath     string
}

func GenerateRegistry(req RegistryRequest) error {

	if !IsExist(fmt.Sprintf("%s/controller/%s/%s.go", req.FolderPath, req.ControllerName, req.UsecaseName)) {
		return fmt.Errorf("controller %s/%s is not found", req.ControllerName, req.UsecaseName)
	}

	if !IsExist(fmt.Sprintf("%s/datasource/%s/%s.go", req.FolderPath, req.DatasourceName, req.UsecaseName)) {
		return fmt.Errorf("datasource %s/%s is not found", req.DatasourceName, req.UsecaseName)
	}

	if !IsExist(fmt.Sprintf("%s/usecase/%s", req.FolderPath, req.UsecaseName)) {
		return fmt.Errorf("usecase %s is not found", req.UsecaseName)
	}

	funcDeclareInjectedCode := fmt.Sprintf(funcDeclare+"\n", CamelCase(req.UsecaseName), req.DatasourceName, req.UsecaseName, LowerCase(req.UsecaseName), req.UsecaseName, LowerCase(req.UsecaseName), req.ControllerName, req.UsecaseName)
	funcCallInjectedCode := fmt.Sprintf(funcCall, CamelCase(req.UsecaseName))

	file, err := os.Open(fmt.Sprintf("%s/application/registry.go", req.FolderPath))
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

	if err := ioutil.WriteFile(fmt.Sprintf("%s/application/registry.go", req.FolderPath), buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

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

func (d *registry) Generate(args ...string) error {

	if len(args) < 5 {
		return fmt.Errorf("please define controller, datasource, usecase name. ex: `gogen registry restapi production CreateMenu`")
	}

	controllerName := args[2]
	datasourceName := args[3]
	usecaseName := args[4]

	if !IsExist(fmt.Sprintf("controller/%s/%s.go", controllerName, usecaseName)) {
		return fmt.Errorf("controller %s/%s is not found", controllerName, usecaseName)
	}

	if !IsExist(fmt.Sprintf("datasource/%s/%s.go", datasourceName, usecaseName)) {
		return fmt.Errorf("datasource %s/%s is not found", datasourceName, usecaseName)
	}

	if !IsExist(fmt.Sprintf("usecase/%s", usecaseName)) {
		return fmt.Errorf("usecase %s is not found", usecaseName)
	}

	funcDeclare := `
func %s(a *Application) {
	outport := %s.New%sDatasource()
	inport := %s.New%sUsecase(outport)
	a.Router.POST("/%s", %s.%s(inport))
}`

	funcDeclareInjectedCode := fmt.Sprintf(funcDeclare+"\n", CamelCase(usecaseName), datasourceName, usecaseName, LowerCase(usecaseName), usecaseName, LowerCase(usecaseName), controllerName, usecaseName)
	funcCallInjectedCode := fmt.Sprintf("	%s(a)", CamelCase(usecaseName))

	file, err := os.Open("application/registry.go")
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

	if err := ioutil.WriteFile("application/registry.go", buffer.Bytes(), 0644); err != nil {
		return err
	}

	return nil
}

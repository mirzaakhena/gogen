package gogen

import (
	"fmt"
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

	output := `
  func %s(a *Application) {
    outport := %s.New%sDatasource()
    inport := %s.New%sUsecase(outport)
    a.Router.POST("/%s", %s.%s(inport))
  }`

	fmt.Printf(output+"\n\n", CamelCase(usecaseName), datasourceName, usecaseName, LowerCase(usecaseName), usecaseName, LowerCase(usecaseName), controllerName, usecaseName)

	return nil
}

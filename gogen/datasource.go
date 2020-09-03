package gogen

import (
	"fmt"
	"strings"
)

const (
	DATASOURCE_NAME_INDEX         int = 2
	DATASOURCE_USECASE_NAME_INDEX int = 3
)

type datasource struct {
}

func NewDatasource() Generator {
	return &datasource{}
}

func (d *datasource) Generate(args ...string) error {

	if IsNotExist(".application_schema/") {
		return fmt.Errorf("please call `gogen init` first")
	}

	if len(args) < 4 {
		return fmt.Errorf("please define datasource and usecase_name. ex: `gogen datasource production CreateOrder`")
	}

	datasourceName := args[DATASOURCE_NAME_INDEX]

	usecaseName := args[DATASOURCE_USECASE_NAME_INDEX]

	if IsNotExist(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName)) {
		return fmt.Errorf("Usecase `%s` is not found. Generate it by call `gogen usecase %s` first", usecaseName, usecaseName)
	}

	tp, err := ReadYAML(usecaseName)
	if err != nil {
		return err
	}

	tp.DatasourceName = datasourceName

	CreateFolder("datasources/%s", strings.ToLower(datasourceName))

	WriteFileIfNotExist(
		"datasources/datasource/datasource._go",
		fmt.Sprintf("datasources/%s/%s.go", datasourceName, usecaseName),
		tp,
	)

	return nil
}

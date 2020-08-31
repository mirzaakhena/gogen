package gogen

import (
	"fmt"
	"os"
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

	{
		_, err := os.Stat(".application_schema/")
		if os.IsNotExist(err) {
			return fmt.Errorf("please call `gogen init` first")
		}
	}

	if len(args) < 4 {
		return fmt.Errorf("please define usecase name. ex: `gogen datasource production CreateOrder`")
	}

	productionName := args[DATASOURCE_NAME_INDEX]

	usecaseName := args[DATASOURCE_USECASE_NAME_INDEX]

	return nil
}

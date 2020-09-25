package gogen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type datasource struct {
}

func NewDatasource() Generator {
	return &datasource{}
}

func (d *datasource) Generate(args ...string) error {

	// if IsNotExist(".application_schema/") {
	// 	return fmt.Errorf("please call `gogen init` first")
	// }

	if len(args) < 4 {
		return fmt.Errorf("please define datasource and usecase_name. ex: `gogen datasource production CreateOrder`")
	}

	datasourceName := args[2]

	usecaseName := args[3]

	// if IsNotExist(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName)) {
	// 	return fmt.Errorf("Usecase `%s` is not found. Generate it by call `gogen usecase %s` first", usecaseName, usecaseName)
	// }

	// tp, err := ReadYAML(usecaseName)
	// if err != nil {
	// 	return err
	// }

	ds := Datasource{}
	ds.DatasourceName = datasourceName
	ds.UsecaseName = usecaseName
	ds.PackagePath = GetPackagePath()

	file, err := os.Open(fmt.Sprintf("usecases/%s/outport/outport.go", usecaseName))
	if err != nil {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	state := 0
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %s interface {", usecaseName)) {
			state = 1
		} else //
		if state == 1 {
			if strings.HasPrefix(scanner.Text(), "}") {
				state = 2
				break
			} else {
				completeMethod := strings.TrimSpace(scanner.Text())
				methodNameOnly := strings.Split(completeMethod, "(")[0]
				fmt.Println(methodNameOnly)
				ds.OutportMethods = append(ds.OutportMethods, methodNameOnly)
			}
		}
	}

	CreateFolder("datasources/%s", strings.ToLower(datasourceName))

	WriteFileIfNotExist(
		"datasources/datasource/datasource._go",
		fmt.Sprintf("datasources/%s/%s.go", datasourceName, usecaseName),
		ds,
	)

	GoFormat(ds.PackagePath)

	return nil
}

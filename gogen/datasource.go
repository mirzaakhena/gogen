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

	if len(args) < 4 {
		return fmt.Errorf("please define datasource and usecase_name. ex: `gogen datasource production CreateOrder`")
	}

	datasourceName := args[2]

	usecaseName := args[3]

	folderPath := "."

	return GenerateDatasource(datasourceName, usecaseName, folderPath)

}

func GenerateDatasource(datasourceName, usecaseName, folderPath string) error {

	var folderImport string
	if folderPath != "." {
		folderImport = fmt.Sprintf("/%s", folderPath)
	}

	ds := Datasource{}
	ds.Directory = folderImport
	ds.DatasourceName = datasourceName
	ds.UsecaseName = usecaseName
	ds.PackagePath = GetPackagePath()

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("error1. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := 0
		for scanner.Scan() {
			if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sOutport interface {", usecaseName)) {
				state = 1
			} else //
			if state == 1 {
				if strings.HasPrefix(scanner.Text(), "}") {
					state = 2
					break
				} else {
					completeMethod := strings.TrimSpace(scanner.Text())
					methodNameOnly := strings.Split(completeMethod, "(")[0]
					ds.Outports = append(ds.Outports, &Outport{
						Name: methodNameOnly,
					})
				}
			}
		}

		if state == 0 {
			return fmt.Errorf("error2. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	for _, ot := range ds.Outports {

		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("error3. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := 0
		for scanner.Scan() {
			if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", ot.Name)) {
				state = 1
			} else //
			if state == 1 {
				if strings.HasPrefix(scanner.Text(), "}") {
					state = 2
					break
				} else {

					completeFieldWithType := strings.TrimSpace(scanner.Text())
					if len(completeFieldWithType) == 0 {
						continue
					}
					fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
					ot.RequestFields = append(ot.RequestFields, &NameType{
						Name: strings.TrimSpace(fieldWithType[0]),
					})

				}
			}
		}

		if state == 0 {
			return fmt.Errorf("error4. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	for _, ot := range ds.Outports {

		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("error5. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := 0
		for scanner.Scan() {
			if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", ot.Name)) {
				state = 1
			} else //
			if state == 1 {
				if strings.HasPrefix(scanner.Text(), "}") {
					state = 2
					break
				} else {

					completeFieldWithType := strings.TrimSpace(scanner.Text())
					if len(completeFieldWithType) == 0 {
						continue
					}
					fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
					ot.ResponseFields = append(ot.ResponseFields, &NameType{
						Name: strings.TrimSpace(fieldWithType[0]),
					})

				}
			}
		}

		if state == 0 {
			return fmt.Errorf("error6. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	CreateFolder("%s/datasource/%s", folderPath, strings.ToLower(datasourceName))

	_ = WriteFileIfNotExist(
		"datasource/datasourceName/datasource._go",
		fmt.Sprintf("%s/datasource/%s/%s.go", folderPath, datasourceName, usecaseName),
		ds,
	)

	GoFormat(ds.PackagePath)

	return nil
}

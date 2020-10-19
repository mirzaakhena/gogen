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

	return GenerateDatasource(DatasourceRequest{
		DatasourceName: args[2],
		UsecaseName:    args[3],
		FolderPath:     ".",
	})

}

type DatasourceRequest struct {
	DatasourceName string
	UsecaseName    string
	FolderPath     string
}

func GenerateDatasource(req DatasourceRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	ds := Datasource{}
	ds.Directory = folderImport
	ds.DatasourceName = req.DatasourceName
	ds.UsecaseName = req.UsecaseName
	ds.PackagePath = GetPackagePath()

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
		if err != nil {
			return fmt.Errorf("error1. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := 0
		for scanner.Scan() {
			if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sOutport interface {", req.UsecaseName)) {
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
			return fmt.Errorf("error2. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
	}

	for _, ot := range ds.Outports {

		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
		if err != nil {
			return fmt.Errorf("error3. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
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
			return fmt.Errorf("error4. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
	}

	for _, ot := range ds.Outports {

		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
		if err != nil {
			return fmt.Errorf("error5. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
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
			return fmt.Errorf("error6. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
	}

	CreateFolder("%s/datasource/%s", req.FolderPath, strings.ToLower(req.DatasourceName))

	_ = WriteFileIfNotExist(
		"datasource/datasourceName/datasource._go",
		fmt.Sprintf("%s/datasource/%s/%s.go", req.FolderPath, req.DatasourceName, req.UsecaseName),
		ds,
	)

	GoFormat(ds.PackagePath)

	return nil
}

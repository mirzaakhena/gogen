package gogen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type test struct {
}

func NewTest() Generator {
	return &test{}
}

func (d *test) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define test usecase. ex: `gogen test CreateOrder`")
	}

	usecaseName := args[2]

	folderPath := "hehe"

	return GenerateTest(usecaseName, folderPath)
}

func GenerateTest(usecaseName, folderPath string) error {

	var folderImport string
	if folderPath != "." {
		folderImport = fmt.Sprintf("/%s", folderPath)
	}

	ds := Test{}
	ds.UsecaseName = usecaseName
	ds.Directory = folderImport
	ds.PackagePath = GetPackagePath()

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := 0
		for scanner.Scan() {
			if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sInport interface {", usecaseName)) {
				state = 1
			} else //
			if state == 1 {
				completeMethod := strings.TrimSpace(scanner.Text())
				methodNameOnly := strings.Split(completeMethod, "(")[0]
				ds.Type = methodNameOnly
				break
			}
		}
		if state == 0 {
			return fmt.Errorf("usecase %s is not found", usecaseName)
		}
	}

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
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
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	for _, ot := range ds.Outports {

		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
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
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	for _, ot := range ds.Outports {

		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
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
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := 0
		for scanner.Scan() {
			if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", usecaseName)) {
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
					ds.InportRequestFields = append(ds.InportRequestFields, &NameType{
						Name: strings.TrimSpace(fieldWithType[0]),
						Type: strings.TrimSpace(fieldWithType[1]),
					})
				}
			}
		}
		if state == 0 {
			return fmt.Errorf("not found Request struct")
		}
	}

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := 0
		for scanner.Scan() {
			if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", usecaseName)) {
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
					ds.InportResponseFields = append(ds.InportResponseFields, &NameType{
						Name: strings.TrimSpace(fieldWithType[0]),
						Type: strings.TrimSpace(fieldWithType[1]),
					})
				}
			}
		}
		if state == 0 {
			return fmt.Errorf("not found Response struct")
		}
	}

	if ds.Type == "HandleQuery" {
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor_test-query._go",
			fmt.Sprintf("%s/usecase/%s/interactor_test.go", folderPath, strings.ToLower(usecaseName)),
			ds,
		)
	} else //

	if ds.Type == "HandleCommand" {
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor_test-command._go",
			fmt.Sprintf("%s/usecase/%s/interactor_test.go", folderPath, strings.ToLower(usecaseName)),
			ds,
		)
	}

	GenerateMock(ds.PackagePath, ds.UsecaseName, folderPath)

	return nil
}

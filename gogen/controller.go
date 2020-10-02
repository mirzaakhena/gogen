package gogen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type controller struct {
}

func NewController() Generator {
	return &controller{}
}

func (d *controller) Generate(args ...string) error {

	if len(args) < 4 {
		return fmt.Errorf("please define datasource and usecase_name. ex: `gogen controller restapi.gin CreateOrder`")
	}

	controllerType := args[2]

	usecaseName := args[3]

	ct := Controller{}
	ct.UsecaseName = usecaseName
	ct.PackagePath = GetPackagePath()

	{
		file, err := os.Open(fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(usecaseName)))
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
				ct.Type = methodNameOnly
				break
			}
		}
		if state == 0 {
			return fmt.Errorf("not found inport method HandleQuery or HandleCommand.")
		}
	}

	{
		file, err := os.Open(fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(usecaseName)))
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
					ct.InportFields = append(ct.InportFields, &NameType{
						Name: strings.TrimSpace(fieldWithType[0]),
						Type: strings.TrimSpace(fieldWithType[1]),
					})
				}
			}
		}
		if state == 0 {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	if controllerType == "restapi.gin" {

		CreateFolder("controller/restapi")

		if ct.Type == "HandleQuery" {
			_ = WriteFileIfNotExist(
				"controller/restapi/gin-query._go",
				fmt.Sprintf("controller/restapi/%s.go", usecaseName),
				ct,
			)
		} else //

		if ct.Type == "HandleCommand" {
			_ = WriteFileIfNotExist(
				"controller/restapi/gin-command._go",
				fmt.Sprintf("controller/restapi/%s.go", usecaseName),
				ct,
			)
		}

	} else //

	if controllerType == "restapi.http" {

		CreateFolder("controller/restapi")

		_ = WriteFileIfNotExist(
			"controller/restapi/http._go",
			fmt.Sprintf("controller/restapi/%s.go", usecaseName),
			ct,
		)

	}

	GoFormat(ct.PackagePath)

	return nil
}

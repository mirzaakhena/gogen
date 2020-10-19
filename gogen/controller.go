package gogen

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/mirzaakhena/templator"
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

	return GenerateController(controllerType, usecaseName)

}

func GenerateController(controllerType, usecaseName string) error {

	ct := Controller{}
	ct.UsecaseName = usecaseName
	ct.PackagePath = templator.GetPackagePath()

	{
		file, err := os.Open(fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(usecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := "FIND_INTERFACE"
		for scanner.Scan() {
			if state == "FIND_INTERFACE" && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sInport interface {", usecaseName)) {
				state = "FIND_METHOD_SIGNATURE"
			} else //
			if state == "FIND_METHOD_SIGNATURE" {
				completeMethod := strings.TrimSpace(scanner.Text())
				methodNameOnly := strings.Split(completeMethod, "(")[0]
				ct.Type = methodNameOnly
				break
			}
		}
		if state == "FIND_INTERFACE" {
			return fmt.Errorf("usecase %s is not found", usecaseName)
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

		state := "FIND_REQUEST_STRUCT"
		for scanner.Scan() {
			if state == "FIND_REQUEST_STRUCT" && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", usecaseName)) {
				state = "FIND_FIELD_AND_TYPE"
			} else //
			if state == "FIND_FIELD_AND_TYPE" {
				if strings.HasPrefix(scanner.Text(), "}") {
					break
				} else //

				{
					completeFieldWithType := strings.TrimSpace(scanner.Text())
					if len(completeFieldWithType) == 0 {
						continue
					}
					fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
					ct.InportRequestFields = append(ct.InportRequestFields, &NameType{
						Name: strings.TrimSpace(fieldWithType[0]),
						Type: strings.TrimSpace(fieldWithType[1]),
					})
				}
			}
		}
		if state == "FIND_REQUEST_STRUCT" {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
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

		state := "FIND_RESPONSE_STRUCT"
		for scanner.Scan() {
			if state == "FIND_RESPONSE_STRUCT" && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", usecaseName)) {
				state = "FIND_FIELD_AND_TYPE"
			} else //
			if state == "FIND_FIELD_AND_TYPE" {
				if strings.HasPrefix(scanner.Text(), "}") {
					break
				} else //

				{
					completeFieldWithType := strings.TrimSpace(scanner.Text())
					if len(completeFieldWithType) == 0 {
						continue
					}
					fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
					ct.InportResponseFields = append(ct.InportResponseFields, &NameType{
						Name: strings.TrimSpace(fieldWithType[0]),
						Type: strings.TrimSpace(fieldWithType[1]),
					})
				}
			}
		}
		if state == "FIND_RESPONSE_STRUCT" {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}
	}

	if controllerType == "restapi.gin" {

		templator.CreateFolder("controller/restapi")

		if ct.Type == "HandleQuery" {
			_ = templator.WriteFileIfNotExist(
				"controller/restapi/gin-query._go",
				fmt.Sprintf("controller/restapi/%s.go", usecaseName),
				ct,
			)
		} else //

		if ct.Type == "HandleCommand" {
			_ = templator.WriteFileIfNotExist(
				"controller/restapi/gin-command._go",
				fmt.Sprintf("controller/restapi/%s.go", usecaseName),
				ct,
			)
		}

	} else //

	if controllerType == "restapi.http" {

		templator.CreateFolder("controller/restapi")

		_ = templator.WriteFileIfNotExist(
			"controller/restapi/http._go",
			fmt.Sprintf("controller/restapi/%s.go", usecaseName),
			ct,
		)

	}

	templator.GoFormat(ct.PackagePath)

	return nil
}

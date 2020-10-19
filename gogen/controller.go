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

	return GenerateController(ControllerRequest{
		ControllerType: args[2],
		UsecaseName:    args[3],
		FolderPath:     ".",
	})

}

type ControllerRequest struct {
	ControllerType string
	UsecaseName    string
	FolderPath     string
}

func GenerateController(req ControllerRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	ct := Controller{}
	ct.UsecaseName = req.UsecaseName
	ct.Directory = folderImport
	ct.PackagePath = GetPackagePath()

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := "FIND_INTERFACE"
		for scanner.Scan() {
			if state == "FIND_INTERFACE" && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sInport interface {", req.UsecaseName)) {
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
			return fmt.Errorf("usecase %s is not found", req.UsecaseName)
		}
	}

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := "FIND_REQUEST_STRUCT"
		for scanner.Scan() {
			if state == "FIND_REQUEST_STRUCT" && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", req.UsecaseName)) {
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
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
	}

	{
		file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
		if err != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		scanner.Split(bufio.ScanLines)

		state := "FIND_RESPONSE_STRUCT"
		for scanner.Scan() {
			if state == "FIND_RESPONSE_STRUCT" && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", req.UsecaseName)) {
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
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}
	}

	if req.ControllerType == "restapi.gin" {

		CreateFolder("%s/controller/restapi", req.FolderPath)

		if ct.Type == "HandleQuery" {
			_ = WriteFileIfNotExist(
				"controller/restapi/gin-query._go",
				fmt.Sprintf("%s/controller/restapi/%s.go", req.FolderPath, req.UsecaseName),
				ct,
			)
		} else //

		if ct.Type == "HandleCommand" {
			_ = WriteFileIfNotExist(
				"controller/restapi/gin-command._go",
				fmt.Sprintf("%s/controller/restapi/%s.go", req.FolderPath, req.UsecaseName),
				ct,
			)
		}

	} else //

	if req.ControllerType == "restapi.http" {

		CreateFolder("%s/controller/restapi", req.FolderPath)

		_ = WriteFileIfNotExist(
			"controller/restapi/http._go",
			fmt.Sprintf("%s/controller/restapi/%s.go", req.FolderPath, req.UsecaseName),
			ct,
		)

	}

	GoFormat(ct.PackagePath)

	return nil
}

package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
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

	inportFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName))
	node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
	if errParse != nil {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
	}

	ct.InportRequestFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ct.UsecaseName, "Request"))

	ct.InportResponseFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ct.UsecaseName, "Response"))

	if req.ControllerType == "restapi.gin" {

		CreateFolder("%s/controller/restapi", req.FolderPath)

		_ = WriteFileIfNotExist(
			"controller/restapi/gin._go",
			fmt.Sprintf("%s/controller/restapi/%s.go", req.FolderPath, req.UsecaseName),
			ct,
		)

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

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
		return fmt.Errorf("please define package_name and usecase_name. ex: `gogen controller restapi CreateOrder`")
	}

	return GenerateController(ControllerRequest{
		ControllerPackage: args[2],
		UsecaseName:       args[3],
		FolderPath:        ".",
	})

}

type ControllerRequest struct {
	ControllerPackage string
	UsecaseName       string
	FolderPath        string
}

func GenerateController(req ControllerRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	uc := Usecase{
		PackageName: req.ControllerPackage,
		Name:        req.UsecaseName,
		Directory:   folderImport,
		PackagePath: GetPackagePath(),
	}

	inportFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName))
	node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
	if errParse != nil {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
	}

	uc.InportRequestFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", req.UsecaseName, "Request"))

	uc.InportResponseFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", req.UsecaseName, "Response"))

	CreateFolder("%s/controller/%s", req.FolderPath, LowerCase(req.ControllerPackage))

	_ = WriteFileIfNotExist(
		"controller/basic._go",
		fmt.Sprintf("%s/controller/%s/%s.go", req.FolderPath, LowerCase(req.ControllerPackage), req.UsecaseName),
		uc,
	)

	// if req.ControllerPackage == "restapi.gin" {

	// 	CreateFolder("%s/controller/restapi", req.FolderPath)

	// 	_ = WriteFileIfNotExist(
	// 		"controller/restapi/gin._go",
	// 		fmt.Sprintf("%s/controller/restapi/%s.go", req.FolderPath, req.UsecaseName),
	// 		uc,
	// 	)

	// } else //

	// if req.ControllerPackage == "restapi.http" {

	// 	CreateFolder("%s/controller/restapi", req.FolderPath)

	// 	_ = WriteFileIfNotExist(
	// 		"controller/restapi/http._go",
	// 		fmt.Sprintf("%s/controller/restapi/%s.go", req.FolderPath, req.UsecaseName),
	// 		uc,
	// 	)

	// }

	GoFormat(uc.PackagePath)

	return nil
}

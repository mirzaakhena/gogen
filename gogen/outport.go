package gogen

import (
	"fmt"
	"strings"
)

type outport struct {
}

func NewOutport() Generator {
	return &outport{}
}

func (d *outport) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen outport CreateOrder CheckOrder SaveOrder`")
	}

	return GenerateOutport(OutportRequest{
		UsecaseName: args[2],
		Methods:     args[3:],
		FolderPath:  ".",
	})
}

type OutportRequest struct {
	UsecaseName string
	Methods     []string
	FolderPath  string
}

func GenerateOutport(req OutportRequest) error {

	if IsExist(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName))) {
		return fmt.Errorf("outport file is already exist")
	}

	if !IsExist(fmt.Sprintf("%s/usecase/%s/", req.FolderPath, strings.ToLower(req.UsecaseName))) {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
	}

	_ = WriteFileIfNotExist(
		"usecase/usecaseName/port/outport._go",
		fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)),
		struct {
			Name    string
			Methods []string
		}{
			Name:    req.UsecaseName,
			Methods: req.Methods,
		},
	)

	// CreateFolder("model/")

	// _ = WriteFileIfNotExist(
	// 	"model/model._go",
	// 	fmt.Sprintf("%s/model/%s.go", req.FolderPath, req.ModelName),
	// 	struct {
	// 		Name      string
	// 		Directory string
	// 	}{
	// 		Name:      req.ModelName,
	// 		Directory: folderImport,
	// 	},
	// )

	return nil
}

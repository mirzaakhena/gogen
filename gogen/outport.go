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
		return fmt.Errorf("please define usecase name. ex: `gogen outports CreateOrder CheckOrder SaveOrder`")
	}

	methods := []OutportStruct{}
	for _, m := range args[3:] {
		methods = append(methods, OutportStruct{Name: m})
	}

	return GenerateOutport(OutportRequest{
		UsecaseName: args[2],
		Outports:    methods,
		FolderPath:  ".",
	})
}

type OutportRequest struct {
	UsecaseName string
	Outports    []OutportStruct
	FolderPath  string
}

type OutportStruct struct {
	Name string
}

func GenerateOutport(req OutportRequest) error {

	if IsExist(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName))) {
		return fmt.Errorf("outport file is already exist")
	}

	if !IsExist(fmt.Sprintf("%s/usecase/%s/", req.FolderPath, strings.ToLower(req.UsecaseName))) {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
	}

	_ = WriteFileIfNotExist(
		"usecase/usecaseName/port/outportx._go",
		fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)),
		struct {
			Name     string
			Outports []OutportStruct
		}{
			Name:     req.UsecaseName,
			Outports: req.Outports,
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

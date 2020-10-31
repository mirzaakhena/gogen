package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
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

	outportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName))
	node, errParse := parser.ParseFile(token.NewFileSet(), outportFile, nil, parser.ParseComments)
	if errParse != nil {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
	}

	interfaceNames, err := ReadInterfaceMethodName(node, fmt.Sprintf("%s%s", req.UsecaseName, "Outport"))
	if err != nil {
		return fmt.Errorf("usecase %s is not found 111", req.UsecaseName)
	}

	outportMethods := []*OutportMethod{}

	// existing method name in outport.go
	methodNameMap := map[string]int{}
	for index, name := range interfaceNames {
		methodNameMap[name] = index
		outportMethods = append(outportMethods, &OutportMethod{
			Name: name,
		})
	}

	// checking new method name
	for _, out := range req.Outports {
		if _, exist := methodNameMap[out.Name]; exist {
			continue
		}
		outportMethods = append(outportMethods, &OutportMethod{
			Name: out.Name,
		})
		methodNameMap[out.Name] = len(interfaceNames) - 1
	}

	outport := Outport{
		UsecaseName: req.UsecaseName,
		Methods:     outportMethods,
	}

	// Create outport file
	_ = WriteFile(
		"usecase/usecaseName/port/outport._go",
		fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)),
		outport,
	)

	return nil
}

package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type OutportBuilderRequest struct {
	UsecaseName        string
	FolderPath         string
	OutportMethodNames []string
}

type outportBuilder struct {
	OutportBuilderRequest OutportBuilderRequest
}

func NewOutport(req OutportBuilderRequest) Generator {
	return &outportBuilder{OutportBuilderRequest: req}
}

func (d *outportBuilder) Generate() error {

	usecaseName := d.OutportBuilderRequest.UsecaseName
	folderPath := d.OutportBuilderRequest.FolderPath
	newOutportMethodNames := d.OutportBuilderRequest.OutportMethodNames

	if len(usecaseName) == 0 {
		return fmt.Errorf("Usecase name must not empty")
	}

	if len(newOutportMethodNames) == 0 {
		return fmt.Errorf("Outport name is not defined. At least have one outport name")
	}

	// check the file outport.go exist or not
	outportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName))
	if !IsExist(outportFile) {
		return fmt.Errorf("outport file is not found")
	}

	// prepare fileSet
	fSet := token.NewFileSet()
	node, errParse := parser.ParseFile(fSet, outportFile, nil, parser.ParseComments)
	if errParse != nil {
		return errParse
	}

	// read all struct under outport.go only
	mapStruct := map[string][]FieldType{}
	var errExist error
	mapStruct, errExist = ReadAllStructInFile(node, mapStruct)
	if errExist != nil {
		return errExist
	}

	// read all method in Outport interface
	oldOutportMethods, errRead := ReadInterfaceMethodAndField(node, fmt.Sprintf("%sOutport", usecaseName), mapStruct)
	if errRead != nil {
		return errRead
	}

	// map all method name
	methodNameMap := map[string]int{}
	for _, ms := range oldOutportMethods {
		methodNameMap[ms.MethodName] = 1
	}

	// checking new method name
	for _, newOutportMethodName := range newOutportMethodNames {

		// if found old method then ignore and continue
		if _, exist := methodNameMap[newOutportMethodName]; exist {
			continue
		}

		// add new method
		oldOutportMethods = append(oldOutportMethods, InterfaceMethod{
			MethodName: newOutportMethodName,
			ParamType:  fmt.Sprintf("%sRequest", newOutportMethodName),
			ResultType: fmt.Sprintf("%sResponse", newOutportMethodName),
		})
		methodNameMap[newOutportMethodName] = 1
	}

	// collect non existing
	so := StructureOutport{
		UsecaseName:    usecaseName,
		Methods:        oldOutportMethods,
		ParamsRequired: true,
	}
	_ = WriteFile(
		"usecase/usecaseName/port/outport._go",
		fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)),
		so,
	)

	return nil
}

package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type UsecaseBuilderRequest struct {
	UsecaseName        string
	FolderPath         string
	GomodPath          string
	OutportMethodNames []string
}

type usecaseBuilder struct {
	UsecaseBuilderRequest UsecaseBuilderRequest
}

func NewUsecase(req UsecaseBuilderRequest) Generator {
	return &usecaseBuilder{UsecaseBuilderRequest: req}
}

func (d *usecaseBuilder) Generate() error {

	usecaseName := strings.TrimSpace(d.UsecaseBuilderRequest.UsecaseName)
	folderPath := d.UsecaseBuilderRequest.FolderPath
	outportMethodNames := d.UsecaseBuilderRequest.OutportMethodNames
	gomodPath := d.UsecaseBuilderRequest.GomodPath

	if len(usecaseName) == 0 {
		return fmt.Errorf("Usecase name must not empty")
	}

	// create a folder with usecase name
	CreateFolder("%s/usecase/%s/port", folderPath, strings.ToLower(usecaseName))

	// create a port/inport.go file
	{
		si := StructureInport{
			UsecaseName: usecaseName,
		}
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/inport._go",
			fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(usecaseName)),
			si,
		)
	}

	// create a port/outport.go file
	{
		outportMethods := []InterfaceMethod{}
		paramsRequired := true
		if len(outportMethodNames) > 0 {
			for _, methodName := range outportMethodNames {
				outportMethods = append(outportMethods, InterfaceMethod{
					MethodName: methodName,
				})
			}
		} else {
			outportMethods = append(outportMethods, InterfaceMethod{
				MethodName: usecaseName,
			})
			paramsRequired = false
		}
		so := StructureOutport{
			UsecaseName:    usecaseName,
			Methods:        outportMethods,
			ParamsRequired: paramsRequired,
		}
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/outport._go",
			fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName)),
			so,
		)
	}

	// create a interactor.go file
	{
		mapStruct, errCollect := CollectPortStructs(folderPath, usecaseName)
		if errCollect != nil {
			return errCollect
		}

		uc, errConstruct := ConstructStructureUsecase(gomodPath, folderPath, usecaseName, mapStruct)
		if errConstruct != nil {
			return errConstruct
		}

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor._go",
			fmt.Sprintf("%s/usecase/%s/interactor.go", folderPath, strings.ToLower(usecaseName)),
			uc,
		)
	}

	return nil
}

func CollectPortStructs(folderPath, usecaseName string) (map[string][]FieldType, error) {

	mapStruct := map[string][]FieldType{}

	files, err := ReadAllFileUnderFolder(fmt.Sprintf("%s/usecase/%s/port", folderPath, strings.ToLower(usecaseName)))
	if err != nil {
		return nil, err
	}

	var errExist error
	for _, filename := range files {

		portFile := fmt.Sprintf("%s/usecase/%s/port/%s", folderPath, strings.ToLower(usecaseName), filename)
		fSet := token.NewFileSet()
		node, errParse := parser.ParseFile(fSet, portFile, nil, parser.ParseComments)
		if errParse != nil {
			return nil, errParse
		}

		mapStruct, errExist = ReadAllStructInFile(node, mapStruct)
		if errExist != nil {
			return nil, errExist
		}

	}

	return mapStruct, nil

}

func ConstructStructureUsecase(gomodPath, folderPath, usecaseName string, mapStruct map[string][]FieldType) (*StructureUsecase, error) {

	inportMethod := InterfaceMethod{}
	var outportMethods []InterfaceMethod

	// INPORT
	{
		inportFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(usecaseName))
		fSet := token.NewFileSet()
		node, errParse := parser.ParseFile(fSet, inportFile, nil, parser.ParseComments)
		if errParse != nil {
			return nil, errParse
		}

		inportMethods, err := ReadInterfaceMethodAndField(node, fmt.Sprintf("%sInport", PascalCase(usecaseName)), mapStruct)
		if err != nil {
			return nil, err
		}

		if len(inportMethods) == 1 {
			inportMethod = inportMethods[0]
		}
	}

	// OUTPORT
	{
		outportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName))
		fSet := token.NewFileSet()
		node, errParse := parser.ParseFile(fSet, outportFile, nil, parser.ParseComments)
		if errParse != nil {
			return nil, errParse
		}

		var errRead error
		outportMethods, errRead = ReadInterfaceMethodAndField(node, fmt.Sprintf("%sOutport", PascalCase(usecaseName)), mapStruct)
		if errRead != nil {
			return nil, errRead
		}

	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	uc := StructureUsecase{
		UsecaseName: usecaseName,
		PackagePath: packagePath,
		Inport:      inportMethod,
		Outport:     outportMethods,
	}

	return &uc, nil

}

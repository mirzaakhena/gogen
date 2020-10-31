package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type usecase struct {
}

func NewUsecase() Generator {
	return &usecase{}
}

func (d *usecase) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
	}

	return GenerateUsecase(UsecaseRequest{
		UsecaseName: args[2],
		FolderPath:  ".",
		// OutportMethods: []string{"DoSomething"},
	})
}

type UsecaseRequest struct {
	UsecaseName    string
	OutportMethods []string
	FolderPath     string
}

func GenerateUsecase(req UsecaseRequest) error {

	packagePath := GetPackagePath()

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	uc := Usecase{
		Name:        req.UsecaseName,
		PackagePath: packagePath,
		Directory:   folderImport,
	}

	firstTime := false

	if !IsExist(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(uc.Name))) {
		if len(req.OutportMethods) == 0 {
			req.OutportMethods = []string{"DoSomething"}
			firstTime = true
		}
	}

	// set outport methods
	{
		outports := []*Outport{}
		for _, methodName := range req.OutportMethods {
			outports = append(outports, &Outport{
				Name:           methodName,
				RequestFields:  nil,
				ResponseFields: nil,
			})
		}
		uc.Outports = outports
	}

	// Create Port Folder
	CreateFolder("%s/usecase/%s/port", req.FolderPath, strings.ToLower(uc.Name))

	// Create inport file
	_ = WriteFileIfNotExist(
		"usecase/usecaseName/port/inport._go",
		fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(uc.Name)),
		uc,
	)

	// Create outport file
	_ = WriteFileIfNotExist(
		"usecase/usecaseName/port/outport._go",
		fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(uc.Name)),
		uc,
	)

	err := readInportOutport(&uc, req.FolderPath, req.UsecaseName, firstTime)
	if err != nil {
		return err
	}

	// create interactor file
	_ = WriteFileIfNotExist(
		"usecase/usecaseName/interactor._go",
		fmt.Sprintf("%s/usecase/%s/interactor.go", req.FolderPath, strings.ToLower(uc.Name)),
		uc,
	)

	return nil
}

func readInportOutport(uc *Usecase, folderPath, usecaseName string, firstTime bool) error {
	{
		inportFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(usecaseName))
		node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
		if errParse != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}

		uc.InportRequestFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", usecaseName, "Request"))
		uc.InportResponseFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", usecaseName, "Response"))
	}

	if !firstTime {

		inportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", folderPath, strings.ToLower(usecaseName))
		node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
		if errParse != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
		}

		interfaceNames, err := ReadInterfaceMethodName(node, fmt.Sprintf("%s%s", usecaseName, "Outport"))
		if err != nil {
			return fmt.Errorf("usecase %s is not found 111", usecaseName)
		}

		for _, methodName := range interfaceNames {
			uc.Outports = append(uc.Outports, &Outport{
				Name: methodName,
			})

			for _, ot := range uc.Outports {
				ot.RequestFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ot.Name, "Request"))
				ot.ResponseFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ot.Name, "Response"))
			}

		}

	}

	return nil
}

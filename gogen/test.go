package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

type test struct {
}

func NewTest() Generator {
	return &test{}
}

func (d *test) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define test usecase. ex: `gogen test CreateOrder`")
	}

	return GenerateTest(TestRequest{
		UsecaseName: args[2],
		FolderPath:  ".",
	})
}

type TestRequest struct {
	UsecaseName string
	FolderPath  string
}

func GenerateTest(req TestRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	ds := Test{}
	ds.UsecaseName = req.UsecaseName
	ds.Directory = folderImport
	ds.PackagePath = GetPackagePath()

	{

		inportFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName))
		node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
		if errParse != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}

		ds.InportRequestFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", req.UsecaseName, "Request"))

		ds.InportResponseFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", req.UsecaseName, "Response"))
	}

	{

		inportFile := fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName))
		node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
		if errParse != nil {
			return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
		}

		interfaceNames, err := ReadInterfaceMethodName(node, fmt.Sprintf("%s%s", req.UsecaseName, "Outport"))
		if err != nil {
			return fmt.Errorf("usecase %s is not found 111", req.UsecaseName)
		}

		for _, methodName := range interfaceNames {
			ds.Outports = append(ds.Outports, &Outport{
				Name: methodName,
			})

			for _, ot := range ds.Outports {
				ot.RequestFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ot.Name, "Request"))
				ot.ResponseFields = ReadFieldInStruct(node, fmt.Sprintf("%s%s", ot.Name, "Response"))
			}

		}

	}

	_ = WriteFileIfNotExist(
		"usecase/usecaseName/interactor_test._go",
		fmt.Sprintf("%s/usecase/%s/interactor_test.go", req.FolderPath, strings.ToLower(req.UsecaseName)),
		ds,
	)

	GenerateMock(ds.PackagePath, ds.UsecaseName, req.FolderPath)

	return nil
}

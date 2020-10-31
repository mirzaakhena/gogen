package gogen

import (
	"fmt"
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

	uc := Usecase{
		Name:        req.UsecaseName,
		Directory:   folderImport,
		PackagePath: GetPackagePath(),
		Outport: &Outport{
			UsecaseName: req.UsecaseName,
		},
	}

	if err := readInport(&uc, req.FolderPath, req.UsecaseName); err != nil {
		return err
	}

	if err := readOutport(uc.Outport, req.FolderPath, req.UsecaseName); err != nil {
		return err
	}

	_ = WriteFileIfNotExist(
		"usecase/usecaseName/interactor_test._go",
		fmt.Sprintf("%s/usecase/%s/interactor_test.go", req.FolderPath, strings.ToLower(req.UsecaseName)),
		uc,
	)

	GenerateMock(uc.PackagePath, uc.Name, req.FolderPath)

	return nil
}

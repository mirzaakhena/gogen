package gogen

import (
	"fmt"
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
	FolderPath     string
	OutportMethods []string
}

func GenerateUsecase(req UsecaseRequest) error {

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	outportMethods := []*OutportMethod{}

	for _, m := range req.OutportMethods {
		outportMethods = append(outportMethods, &OutportMethod{
			Name: m,
		})
	}

	uc := Usecase{
		Name:        req.UsecaseName,
		Directory:   folderImport,
		PackagePath: GetPackagePath(),
		Outport: &Outport{
			UsecaseName: req.UsecaseName,
			Methods:     outportMethods,
		},
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
		uc.Outport,
	)

	if err := readInport(&uc, req.FolderPath, req.UsecaseName); err != nil {
		return err
	}

	if err := readOutport(uc.Outport, req.FolderPath, req.UsecaseName); err != nil {
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

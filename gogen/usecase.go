package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"io/ioutil"
	"log"
	"strings"
)

type usecase struct {
}

func NewUsecase() Generator {
	return &usecase{}
}

func (d *usecase) Generate(args ...string) error {

	if len(args) == 2 {
		DisplayUsecase(".")
		return nil
		//fmt.Errorf("please define usecase name. ex: `gogen usecase CreateOrder`")
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

func DisplayUsecase(folderPath string) {

	files, errReadDir := ioutil.ReadDir(fmt.Sprintf("%s/usecase/", folderPath))
	if errReadDir != nil {
		log.Fatal(errReadDir)
	}

	interfaceNames := []string{}
	for _, file := range files {

		if IsExist(fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, file.Name())) {

			inportFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, file.Name())
			node, errParse := parser.ParseFile(token.NewFileSet(), inportFile, nil, parser.ParseComments)
			if errParse != nil {
				continue
			}

			names := ReadInterfaceName(node)
			for _, intfName := range names {
				if strings.ToLower(intfName) == fmt.Sprintf("%sinport", file.Name()) {
					interfaceNames = append(interfaceNames, string([]rune(intfName))[0:len(intfName)-6])
				}
			}

		}

	}

	if len(interfaceNames) > 0 {
		fmt.Printf("You have %d usecase\n", len(interfaceNames))
		for _, name := range interfaceNames {
			fmt.Printf(" - %s\n", name)
		}
	} else {
		fmt.Printf("You don't have any usecase yet\n")
	}
	fmt.Printf("\nCreate new usecase by this command : `gogen usecase your_usecase_name_in_pascal_cases`\n")

}

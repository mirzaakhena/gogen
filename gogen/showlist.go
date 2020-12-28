package gogen

import (
	"fmt"
	"go/parser"
	"go/token"
	"strings"
)

func ShowAllDatasource() {

}

func ShowAllUsecase() {

	folderPath := "."

	folders, err := ReadAllFileUnderFolder("usecase")
	if err != nil {
		panic(err)
	}

	for _, folderName := range folders {
		portFile := fmt.Sprintf("%s/usecase/%s/port/inport.go", folderPath, strings.ToLower(folderName))
		fSet := token.NewFileSet()
		node, errParse := parser.ParseFile(fSet, portFile, nil, parser.ParseComments)
		if errParse != nil {
			panic(errParse)
		}

		names, errRead := ReadInterfaceName(node)
		if errRead != nil {
			panic(errRead)
		}

		if len(names) > 0 && strings.HasSuffix(names[0], "Inport") {
			fmt.Printf("%s\n", names[0][:len(names[0])-6])
		}

	}

}

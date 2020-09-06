package gogen

import (
	"bufio"
	"fmt"
	"os"
)

const (
	BINDER_USECASE_NAME_INDEX int = 2
)

type binder struct {
}

func NewBinder() Generator {
	return &binder{}
}

func (d *binder) Generate(args ...string) error {

	if IsNotExist(".application_schema/") {
		return fmt.Errorf("please call `gogen init` first")
	}

	if len(args) < 3 {
		return fmt.Errorf("please define usecase_name. ex: `gogen bind CreateOrder`")
	}

	usecaseName := args[BINDER_USECASE_NAME_INDEX]

	if IsNotExist(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName)) {
		return fmt.Errorf("Usecase `%s` is not found. Generate it by call `gogen usecase %s` first", usecaseName, usecaseName)
	}

	if !UsecaseIsExist(usecaseName) {
		InjectCode(usecaseName)
	}

	return nil
}

func UsecaseIsExist(usecaseName string) bool {

	file, err := os.Open(fmt.Sprintf("%s/src/%s/binder/wiring_component.go", GetGopath(), GetPackagePath()))
	if err != nil {
		return false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		row := scanner.Text()
		// 	buffer.WriteString(row)
		// 	buffer.WriteString("\n")
		fmt.Println(row)
	}

	return false
}

func InjectCode(usecaseName string) {

}

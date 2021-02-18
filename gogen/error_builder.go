package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type ErrorBuilderRequest struct {
	ErrorName  string
	FolderPath string
	GomodPath  string
}

type errorBuilder struct {
	ErrorBuilderRequest ErrorBuilderRequest
}

func NewError(req ErrorBuilderRequest) Generator {
	return &errorBuilder{ErrorBuilderRequest: req}
}

func (d *errorBuilder) Generate() error {

	errorName := strings.TrimSpace(d.ErrorBuilderRequest.ErrorName)
	folderPath := d.ErrorBuilderRequest.FolderPath
	gomodPath := d.ErrorBuilderRequest.GomodPath

	if len(errorName) == 0 {
		return fmt.Errorf("ErrorName name must not empty")
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	errorPrefix := GetDefaultErrorEnumPrefix()

	en := StructureError{
		PackagePath: packagePath,
		ErrorName:   errorName,
		ErrorPrefix: errorPrefix,
	}

	CreateFolder("%s/shared/errcat", folderPath)

	_ = WriteFileIfNotExist(
		"shared/errcat/error._go",
		fmt.Sprintf("%s/shared/errcat/error.go", folderPath),
		en,
	)

	_ = WriteFileIfNotExist(
		"shared/errcat/error_enum._go",
		fmt.Sprintf("%s/shared/errcat/error_enum.go", folderPath),
		struct{}{},
	)

	errorFile := fmt.Sprintf("%s/shared/errcat/error_enum.go", folderPath)
	file, err := os.Open(errorFile)
	if err != nil {
		return fmt.Errorf("not found error file")
	}
	defer file.Close()

	constTemplateCode, err := PrintTemplate("shared/errcat/error_enum_template._go", en)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	var constFound bool = false
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		if constFound && strings.HasPrefix(row, ")") {
			constFound = false
			buffer.WriteString(constTemplateCode)
			buffer.WriteString("\n")

		} else //

		if strings.HasPrefix(row, "const (") {
			constFound = true

		}

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	if err := ioutil.WriteFile(fmt.Sprintf("%s/shared/errcat/error_enum.go", folderPath), buffer.Bytes(), 0644); err != nil {
		return err
	}

	GoFormat(packagePath)

	return nil
}

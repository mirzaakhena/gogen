package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type RepositoryBuilderRequest struct {
	RepositoryName string
	EntityName     string
	FolderPath     string
	GomodPath      string
}

type repositoryBuilder struct {
	RepositoryBuilderRequest RepositoryBuilderRequest
}

func NewRepository(req RepositoryBuilderRequest) Generator {
	return &repositoryBuilder{RepositoryBuilderRequest: req}
}

func (d *repositoryBuilder) Generate() error {

	repositoryName := strings.TrimSpace(d.RepositoryBuilderRequest.RepositoryName)
	entityName := strings.TrimSpace(d.RepositoryBuilderRequest.EntityName)
	folderPath := d.RepositoryBuilderRequest.FolderPath
	gomodPath := d.RepositoryBuilderRequest.GomodPath

	if len(repositoryName) == 0 {
		return fmt.Errorf("RepositoryName name must not empty")
	}

	if len(entityName) == 0 {
		return fmt.Errorf("EntityName name must not empty")
	}

	packagePath := GetPackagePath()

	if len(strings.TrimSpace(packagePath)) == 0 {
		packagePath = gomodPath
	}

	en := StructureRepository{
		PackagePath:    packagePath,
		RepositoryName: repositoryName,
		EntityName:     entityName,
	}

	CreateFolder("%s/domain/repository", folderPath)

	_ = WriteFileIfNotExist(
		"domain/repository/repository._go",
		fmt.Sprintf("%s/domain/repository/repository._go", folderPath),
		struct{}{},
	)

	_ = WriteFileIfNotExist(
		"domain/repository/database._go",
		fmt.Sprintf("%s/domain/repository/database._go", folderPath),
		struct{}{},
	)

	errorFile := fmt.Sprintf("%s/domain/repository/repository._go", folderPath)
	file, err := os.Open(errorFile)
	if err != nil {
		return fmt.Errorf("not found error file")
	}
	defer file.Close()

	constTemplateCode, err := PrintTemplate("domain/repository/repository_interface._go", en)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	var buffer bytes.Buffer
	for scanner.Scan() {
		row := scanner.Text()

		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	buffer.WriteString(constTemplateCode)
	buffer.WriteString("\n")

	if err := ioutil.WriteFile(fmt.Sprintf("%s/domain/repository/repository._go", folderPath), buffer.Bytes(), 0644); err != nil {
		return err
	}

	GoImport(fmt.Sprintf("%s/domain/repository/repository._go", folderPath))

	return nil
}

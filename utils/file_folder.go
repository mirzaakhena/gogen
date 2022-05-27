package utils

import (
	"os"
	"text/template"
)

func IsFileExist(filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

// CreateFolderIfNotExist ...
func CreateFolderIfNotExist(folderPath string) (bool, error) {
	if IsFileExist(folderPath) {
		return true, nil
	}
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return false, err
	}
	return false, nil
}

// WriteFileIfNotExist ...
func WriteFileIfNotExist(templateFile, outputFilePath string, obj interface{}) (bool, error) {
	if IsFileExist(outputFilePath) {
		return true, nil
	}
	return false, WriteFile(templateFile, outputFilePath, obj)
}

// WriteFile ...
func WriteFile(templateData, outputFilePath string, obj interface{}) error {

	// TODO move it outside
	fileOut, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}

	tpl, err := template.
		New("something").
		Funcs(FuncMap).
		Parse(templateData)

	if err != nil {
		return err
	}

	err = tpl.Execute(fileOut, obj)
	if err != nil {
		return err
	}

	return nil

}

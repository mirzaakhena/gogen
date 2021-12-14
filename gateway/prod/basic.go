package prod

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"golang.org/x/tools/imports"
	"io/ioutil"
	"os"
	"strings"
	"text/template"
)

type basicUtilityGateway struct {
}

// CreateFolderIfNotExist ...
func (r *basicUtilityGateway) CreateFolderIfNotExist(ctx context.Context, folderPath string) (bool, error) {
	if r.IsFileExist(ctx, folderPath) {
		return true, nil
	}
	if err := os.MkdirAll(folderPath, 0755); err != nil {
		return false, err
	}
	return false, nil
}

// WriteFileIfNotExist ...
func (r *basicUtilityGateway) WriteFileIfNotExist(ctx context.Context, templateFile, outputFilePath string, obj interface{}) (bool, error) {
	if r.IsFileExist(ctx, outputFilePath) {
		return true, nil
	}
	return false, r.WriteFile(ctx, templateFile, outputFilePath, obj)
}

// WriteFile ...
func (r *basicUtilityGateway) WriteFile(ctx context.Context, templateData, outputFilePath string, obj interface{}) error {

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

// PrintTemplate ...
func (r *basicUtilityGateway) PrintTemplate(ctx context.Context, templateString string, x interface{}) (string, error) {

	tpl, err := template.New("something").Funcs(FuncMap).Parse(templateString)
	if err != nil {
		return "", err
	}

	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, x); err != nil {
		return "", err
	}

	return buffer.String(), nil

}

// IsFileExist ...
func (r *basicUtilityGateway) IsFileExist(ctx context.Context, filepath string) bool {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return false
	}
	return true
}

// Reformat ...
func (r *basicUtilityGateway) Reformat(ctx context.Context, goFilename string, bytes []byte) error {

	// reformat the import
	newBytes, err := imports.Process(goFilename, bytes, nil)
	if err != nil {
		return err
	}

	// rewrite it
	if err := ioutil.WriteFile(goFilename, newBytes, 0644); err != nil {
		return err
	}

	return nil
}

// GetPackagePath ...
func (r *basicUtilityGateway) GetPackagePath(ctx context.Context) string {

	var gomodPath string

	file, err := os.Open("go.mod")
	if err != nil {
		fmt.Printf("go.mod is not found. Please create it with command `go mod init your/path/project`\n")
		os.Exit(1)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		if strings.HasPrefix(row, "module") {
			moduleRow := strings.Split(row, " ")
			if len(moduleRow) > 1 {
				gomodPath = moduleRow[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	return strings.Trim(gomodPath, "\"")

}

//// WriteFileIfNotExist2 ...
//func (r *basicUtilityGateway) WriteFileIfNotExist2(ctx context.Context, templateInBytes []byte, outputFilePath string, obj interface{}) (bool, error) {
//  if r.IsFileExist(ctx, outputFilePath) {
//    return true, nil
//  }
//
//  fileOut, err := os.Create(outputFilePath)
//  if err != nil {
//    return false, err
//  }
//
//  tpl, err := template.
//    New("something").
//    Funcs(FuncMap).
//    Parse(string(templateInBytes))
//
//  if err != nil {
//    return false, err
//  }
//
//  err = tpl.Execute(fileOut, obj)
//  if err != nil {
//    return false, err
//  }
//
//  return true, nil
//}
//
//// WriteFile2 ...
//func (r *basicUtilityGateway) WriteFile2(ctx context.Context, templateFile, outputFilePath string, obj interface{}) error {
//  var buffer bytes.Buffer
//
//  scanner := bufio.NewScanner(bytes.NewReader([]byte(templateFile)))
//
//  for scanner.Scan() {
//    row := scanner.Text()
//    buffer.WriteString(row)
//    buffer.WriteString("\n")
//  }
//
//  tpl := template.Must(template.New("something").Funcs(FuncMap).Parse(buffer.String()))
//
//  fileOut, err := os.Create(outputFilePath)
//  if err != nil {
//    return err
//  }
//
//  if err := tpl.Execute(fileOut, obj); err != nil {
//    return err
//  }
//
//  return nil
//}

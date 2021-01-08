package gogen

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"go/build"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

type Generator interface {
	Generate() error
}

func GetGopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func GetPackagePath() string {
	currentPath, _ := filepath.Abs("./")
	pathAfterGopath := strings.Split(currentPath, GetGopath()+"/src/")
	return pathAfterGopath[1]
}

func WriteFileIfNotExist(templateFile, outputFile string, data interface{}) error {

	if !IsExist(outputFile) {
		return WriteFile(templateFile, outputFile, data)
	}

	return nil
}

// function that used in templates
var FuncMap = template.FuncMap{
	"CamelCase":  CamelCase,
	"PascalCase": PascalCase,
	"SnakeCase":  SnakeCase,
	"UpperCase":  UpperCase,
	"LowerCase":  LowerCase,
}

func PrintTemplate(templateFile string, x interface{}) (string, error) {

	file := fmt.Sprintf("%s/src/github.com/mirzaakhena/gogen/templates/%s", GetGopath(), templateFile)
	ts := strings.Split(templateFile, "/")
	tpl, errTemplate := template.New(ts[len(ts)-1]).Funcs(FuncMap).ParseFiles(file)
	if errTemplate != nil {
		return "", errTemplate
	}

	var buffer bytes.Buffer
	if err := tpl.Execute(&buffer, x); err != nil {
		return "", err
	}

	return buffer.String(), nil

}

func WriteFile(templateFile, outputFile string, data interface{}) error {

	var buffer bytes.Buffer

	// this process can be refactor later
	{
		// open template file
		file, err := os.Open(fmt.Sprintf("%s/src/github.com/mirzaakhena/gogen/templates/%s", GetGopath(), templateFile))
		if err != nil {
			return err
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)

		for scanner.Scan() {
			row := scanner.Text()
			buffer.WriteString(row)
			buffer.WriteString("\n")
		}
	}

	tpl := template.Must(template.New("something").Funcs(FuncMap).Parse(buffer.String()))

	fileOut, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	{
		err := tpl.Execute(fileOut, data)
		if err != nil {
			return err
		}
	}

	return nil

}

// CamelCase is
func CamelCase(name string) string {

	// hardcoded is bad
	// But we can figure out later
	{
		if name == "IPAddress" {
			return "ipAddress"
		}

		if name == "ID" {
			return "id"
		}
	}

	out := []rune(name)
	out[0] = unicode.ToLower([]rune(name)[0])
	return string(out)
}

// UpperCase is
func UpperCase(name string) string {
	return strings.ToUpper(name)
}

// LowerCase is
func LowerCase(name string) string {
	return strings.ToLower(name)
}

// PascalCase is
func PascalCase(name string) string {
	rs := []rune(name)
	return strings.ToUpper(string(rs[0])) + string(rs[1:])
}

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase is
func SnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func CreateFolder(format string, a ...interface{}) {
	folderName := fmt.Sprintf(format, a...)

	if err := os.MkdirAll(folderName, 0755); err != nil {
		panic(err)
	}
	fmt.Printf("%v\n", folderName)
}

func GenerateMock(packagePath, usecaseName, folderPath string) {
	fmt.Printf("mockery %s\n", usecaseName)

	lowercaseUsecaseName := strings.ToLower(usecaseName)

	cmd := exec.Command(
		"mockery",

		// read file under path usecase/USECASE_NAME/port/
		"--dir", fmt.Sprintf("%s/usecase/%s/port/", folderPath, lowercaseUsecaseName),

		// specifically interface with name that has suffix 'Outport'
		"--name", fmt.Sprintf("%sOutport", usecaseName),

		// put the mock under usecase/%s/mocks/
		"-output", fmt.Sprintf("%s/usecase/%s/mocks/", folderPath, lowercaseUsecaseName),
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("mockery failed with %s\n", err)
	}
}

func IsExist(fileOrDir string) bool {
	_, err := os.Stat(fileOrDir)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false

}

func GoFormat(path string) {
	fmt.Println("go fmt")
	cmd := exec.Command("go", "fmt", fmt.Sprintf("%s/...", path))
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func ReadAllFileUnderFolder(folderPath string) ([]string, error) {
	var files []string
	fileInfo, err := ioutil.ReadDir(folderPath)
	if err != nil {
		return nil, err
	}

	for _, file := range fileInfo {
		files = append(files, file.Name())
	}

	return files, nil
}

func PrintJSON(x interface{}) string {
	bytes, _ := json.Marshal(x)
	return string(bytes)
}

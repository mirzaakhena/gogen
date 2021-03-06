package util

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
	if len(pathAfterGopath) < 2 {
		return ""
	}
	return pathAfterGopath[1]
}

func WriteFile(templateFile, outputFile string, data interface{}) error {

	var buffer bytes.Buffer

	scanner := bufio.NewScanner(bytes.NewReader([]byte(templateFile)))

	for scanner.Scan() {
		row := scanner.Text()
		buffer.WriteString(row)
		buffer.WriteString("\n")
	}

	tpl := template.Must(template.New("something").Funcs(FuncMap).Parse(buffer.String()))

	fileOut, err := os.Create(outputFile)
	if err != nil {
		return err
	}

	if err := tpl.Execute(fileOut, data); err != nil {
		return err
	}

	return nil

}

func WriteFileIfNotExist(templateFile, outputFile string, data interface{}) error {
	if IsExist(outputFile) {
		return nil
	}
	return WriteFile(templateFile, outputFile, data)
}

// function that used in templates
var FuncMap = template.FuncMap{
	"CamelCase":  CamelCase,
	"PascalCase": PascalCase,
	"SnakeCase":  SnakeCase,
	"UpperCase":  UpperCase,
	"LowerCase":  LowerCase,
	"SpaceCase":  SpaceCase,
}

func PrintTemplate(templateString string, x interface{}) (string, error) {

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

type configStruct struct {
	SelectedTemplate string `json:"template"`
	ErrorPrefix      string `json:"errorPrefix"`
}

func DefaultTemplatePath(templateFile string) string {

	// read local template
	configData, err := ioutil.ReadFile("./.gogen/config.json")
	if err == nil {
		var cs configStruct
		if err := json.Unmarshal(configData, &cs); err != nil {
			panic(err)
		}
		if IsExist(fmt.Sprintf("./.gogen/templates/%s", cs.SelectedTemplate)) {
			return fmt.Sprintf("./.gogen/templates/%s/%s", cs.SelectedTemplate, templateFile)
		}
	}

	// use global default template
	return fmt.Sprintf("%s/src/github.com/mirzaakhena/gogen/templates/default/%s", GetGopath(), templateFile)
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

var matchFirstCapSpaceCase = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCapSpaceCase = regexp.MustCompile("([a-z0-9])([A-Z])")

// SpaceCase is
func SpaceCase(str string) string {
	snake := matchFirstCapSpaceCase.ReplaceAllString(str, "${1} ${2}")
	snake = matchAllCapSpaceCase.ReplaceAllString(snake, "${1} ${2}")
	return strings.ToLower(snake)
}

// PascalCase is
func PascalCase(name string) string {
	rs := []rune(name)
	return strings.ToUpper(string(rs[0])) + string(rs[1:])
}

var matchFirstCapSnakeCase = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCapSnakeCase = regexp.MustCompile("([a-z0-9])([A-Z])")

// SnakeCase is
func SnakeCase(str string) string {
	snake := matchFirstCapSnakeCase.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCapSnakeCase.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func CreateFolderIfNotExist(folderPathWithFormat string, data ...interface{}) error {

	folderName := fmt.Sprintf(folderPathWithFormat, data...)

	if IsExist(folderName) {
		return nil
	}

	if err := os.MkdirAll(folderName, 0755); err != nil {
		return err
	}

	return nil
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

	// cmd.Stdout = os.Stdout
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

func GoFormat() {
	// fmt.Println("go fmt")
	cmd := exec.Command("go", "fmt", fmt.Sprintf("%s/...", GetGoMod()))
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func GoModTidy() {
	// fmt.Println("go mod tidy")
	cmd := exec.Command("go", "mod", "tidy")
	// cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func GoImport(path string) {
	//goimports -w
	cmd := exec.Command("goimports", "-w", path)
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

func GetGoMod() string {

	gomodPath := ""

	file, err := os.Open("go.mod")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

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
		log.Fatal(err)
	}

	return strings.Trim(gomodPath, "\"")

}

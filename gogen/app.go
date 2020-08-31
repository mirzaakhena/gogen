package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/build"
	"html/template"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"unicode"
)

type Generator interface {
	Generate(args ...string) error
}

func getGopath() string {
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		gopath = build.Default.GOPATH
	}
	return gopath
}

func GetPackagePath() string {
	a, _ := filepath.Abs("./")
	x := strings.Split(a, getGopath()+"/src/")
	return x[1]
}

func WriteFileIfNotExist(templateFile, outputFile string, data interface{}) error {

	if IsNotExist(outputFile) {
		return WriteFile(templateFile, outputFile, data)
	}

	return nil
}

func WriteFile(templateFile, outputFile string, data interface{}) error {

	var buffer bytes.Buffer

	// this process can be refactor later
	{

		// open template file
		file, err := os.Open(getGopath() + "/src/github.com/mirzaakhena/gogen/templates/" + templateFile)
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

	// function that used in templates
	funcMap := template.FuncMap{
		"CamelCase":  CamelCase,
		"PascalCase": PascalCase,
		"SnakeCase":  SnakeCase,
		"UpperCase":  UpperCase,
		"LowerCase":  LowerCase,
	}

	tpl := template.Must(template.New("something").Funcs(funcMap).Parse(buffer.String()))

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

	// force it!
	// this is bad. But we can figure out later
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
	return name
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
}

func GenerateMock(packagePath, usecaseName string) {
	fmt.Printf("mockery %s", usecaseName)

	cmd := exec.Command(
		"mockery", "-all",
		// "-case", "snake",
		"-output", "datasources/mocks/",
		"-dir", "outport/")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("cmd.Run() failed with %s\n", err)
	}
}

func IsNotExist(fileOrDir string) bool {
	_, err := os.Stat(fileOrDir)
	return os.IsNotExist(err)
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

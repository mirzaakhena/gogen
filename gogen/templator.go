package gogen

import (
	"bufio"
	"bytes"
	"fmt"
	"go/ast"
	"go/build"
	"go/token"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"text/template"
	"unicode"
)

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

// func ReadYAML(usecaseName string) (*Usecase, error) {

// 	content, err := ioutil.ReadFile(fmt.Sprintf(".application_schema/usecases/%s.yml", usecaseName))
// 	if err != nil {
// 		log.Fatal(err)
// 		return nil, fmt.Errorf("cannot read %s.yml", usecaseName)
// 	}

// 	tp := Usecase{}

// 	if err = yaml.Unmarshal(content, &tp); err != nil {
// 		log.Fatalf("error: %+v", err)
// 		return nil, fmt.Errorf("%s.yml is unrecognized usecase file", usecaseName)
// 	}

// 	tp.Name = usecaseName
// 	tp.PackagePath = GetPackagePath()

// 	return &tp, nil

// }

// func ReadLineByLine(filepath string) []string {
// 	var lineOfCodes []string
// 	{
// 		file, err := os.Open(filepath)
// 		if err != nil {
// 			log.Fatal(err)
// 		}
// 		defer file.Close()

// 		scanner := bufio.NewScanner(file)

// 		scanner.Split(bufio.ScanLines)

// 		for scanner.Scan() {
// 			lineOfCodes = append(lineOfCodes, scanner.Text())
// 		}
// 	}

// 	return lineOfCodes
// }

func ReadInterfaceMethodName(node *ast.File, interfaceName string) ([]string, error) {

	for _, dec := range node.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			if gen.Tok != token.TYPE {
				continue
			}
			for _, specs := range gen.Specs {
				if ts, ok := specs.(*ast.TypeSpec); ok {
					if iface, ok := ts.Type.(*ast.InterfaceType); ok {
						if ts.Name.String() != interfaceName {
							continue
						}
						methods := []string{}
						for _, meths := range iface.Methods.List {
							for _, name := range meths.Names {
								methods = append(methods, name.String())
							}
						}
						return methods, nil

					}
				}
			}
		}
	}
	return nil, fmt.Errorf("interface %s not found", interfaceName)
}

func ReadFieldInStruct(node *ast.File, structName string) []NameType {

	for _, dec := range node.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			if gen.Tok != token.TYPE {
				continue
			}
			for _, specs := range gen.Specs {
				if ts, ok := specs.(*ast.TypeSpec); ok {
					if istruct, ok := ts.Type.(*ast.StructType); ok {
						if !strings.HasSuffix(ts.Name.String(), structName) {
							continue
						}
						nameTypes := []NameType{}
						for _, field := range istruct.Fields.List {
							for _, fieldName := range field.Names {
								nameTypes = append(nameTypes, NameType{Name: fieldName.Name, Type: appendType(field.Type)})
							}
							// ast.Print(fset, field)
						}
						return nameTypes
					}
				}
			}
		}
	}
	return nil
}

func appendType(expr ast.Expr) string {
	var param bytes.Buffer

	for {
		switch t := expr.(type) {
		case *ast.Ident:
			return processIdent(&param, t)

		case *ast.ArrayType:
			return processArrayType(&param, t)

		case *ast.StarExpr:
			return processStarExpr(&param, t)

		case *ast.SelectorExpr:
			return processSelectorExpr(&param, t)

		case *ast.InterfaceType:
			return processInterfaceType(&param, t)

		case *ast.ChanType:
			return processChanType(&param, t)

		case *ast.MapType:
			return processMapType(&param, t)

		case *ast.FuncType:
			return processFuncType(&param, t)

		default:
			return param.String()
		}

	}
}

func processIdent(param *bytes.Buffer, t *ast.Ident) string {
	param.WriteString(t.Name)
	return param.String()
}

func processArrayType(param *bytes.Buffer, t *ast.ArrayType) string {
	if t.Len != nil {
		arrayCapacity := t.Len.(*ast.BasicLit).Value
		param.WriteString(fmt.Sprintf("[%s]", arrayCapacity))
	} else {
		param.WriteString("[]")
	}
	param.WriteString(appendType(t.Elt))

	return param.String()
}

func processStarExpr(param *bytes.Buffer, t *ast.StarExpr) string {
	param.WriteString("*")
	param.WriteString(appendType(t.X))
	return param.String()
}

func processInterfaceType(param *bytes.Buffer, t *ast.InterfaceType) string {
	param.WriteString("interface{}")
	return param.String()
}

func processSelectorExpr(param *bytes.Buffer, t *ast.SelectorExpr) string {
	param.WriteString(appendType(t.X))
	param.WriteString(".")
	param.WriteString(t.Sel.Name)
	return param.String()
}

func processChanType(param *bytes.Buffer, t *ast.ChanType) string {
	if t.Dir == 1 {
		param.WriteString("chan<- ")
	} else if t.Dir == 2 {
		param.WriteString("<-chan ")
	} else {
		param.WriteString("chan ")
	}

	param.WriteString(appendType(t.Value))
	return param.String()
}

func processMapType(param *bytes.Buffer, t *ast.MapType) string {
	param.WriteString("map[")
	param.WriteString(appendType(t.Key))
	param.WriteString("]")
	param.WriteString(appendType(t.Value))
	return param.String()
}

func processFuncType(param *bytes.Buffer, t *ast.FuncType) string {
	param.WriteString("func(")

	nParam := t.Params.NumFields()
	for i, field := range t.Params.List {
		param.WriteString(appendType(field.Type))
		if i+1 < nParam {
			param.WriteString(", ")
		} else {
			param.WriteString(") ")
		}
	}

	nResult := t.Results.NumFields()

	if nResult > 1 {
		param.WriteString("(")
	}

	for i, field := range t.Results.List {
		param.WriteString(appendType(field.Type))
		if i+1 < nResult {
			param.WriteString(", ")
		}
	}

	if nResult > 1 {
		param.WriteString(")")
	}

	return param.String()
}

package gendomain

import (
	"bufio"
	"fmt"
	"github.com/mirzaakhena/gogen/utils"
	"math/rand"
	"os"
	"strings"
	"time"
)

// ObjTemplate ...
type ObjTemplate struct {
	ExecutableName string
	PackagePath    string
	GomodPath      string
	DefaultDomain  string
	DomainName     string
	SecretKey      string
}

func Run(inputs ...string) error {

	if len(inputs) < 1 {
		err := fmt.Errorf("\n" +
			"   # Initiate gogen project with default input. You may change later under .gogen folder\n" +
			"   gogen domain mydomain\n" +
			"     'mydomain' is a your domain name\n" +
			"\n")

		return err
	}

	domainName := inputs[0]

	gomodPath := "your/path/project"
	defaultDomain := fmt.Sprintf("-%s", utils.LowerCase(domainName))

	var letters = []rune("abcdef1234567890")

	randSeq := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letters[rand.Intn(len(letters))]
		}
		return string(b)
	}

	rand.Seed(time.Now().UnixNano())

	obj := &ObjTemplate{
		ExecutableName: utils.GetExecutableName(),
		PackagePath:    utils.GetPackagePath(),
		GomodPath:      gomodPath,
		DefaultDomain:  defaultDomain,
		DomainName:     domainName,
		SecretKey:      randSeq(128),
	}

	fileRenamer := map[string]string{
		"domainname": utils.LowerCase(domainName),
	}

	err := utils.CreateEverythingExactly("templates/", "domain", fileRenamer, obj, utils.AppTemplates)
	if err != nil {
		return err
	}

	_, err = utils.CreateFolderIfNotExist(".gogen")
	if err != nil {
		return err
	}

	gogenDomainFile := "./.gogen/domain"
	exist, err := utils.WriteFileIfNotExist(defaultDomain, gogenDomainFile, struct{}{})
	if err != nil {
		return err
	}

	if exist {
		_ = insertNewDomainName(gogenDomainFile, domainName)
	}

	gitignoreContent := `.idea/
.DS_Store
config.json
*.app
*.exe
*.log
*.db
*/node_modules/
`
	_, err = utils.WriteFileIfNotExist(gitignoreContent, "./.gitignore", struct{}{})
	if err != nil {
		return err
	}

	inFile, err := os.Open(".gogen/domain")
	if err != nil {
		return err
	}
	defer func(inFile *os.File) {
		err := inFile.Close()
		if err != nil {

		}
	}(inFile)

	scanner := bufio.NewScanner(inFile)
	for scanner.Scan() {
		domainNameInGogenFile := strings.TrimSpace(scanner.Text())
		if domainNameInGogenFile == "" {
			continue
		}
		if strings.HasPrefix(domainNameInGogenFile, "-") {
			domainNameInGogenFile = strings.ReplaceAll(domainNameInGogenFile, "-", "")
		}
		domainNameInGogenFile = strings.ToLower(domainNameInGogenFile)
		//_, err := utils.CreateFolderIfNotExist(fmt.Sprintf("domain_%s", domainNameInGogenFile))
		//if err != nil {
		//	return err
		//}

		fileRenamer := map[string]string{
			"domainname": utils.LowerCase(domainNameInGogenFile),
		}

		domainObj := ObjTemplate{DomainName: domainNameInGogenFile}

		err = utils.CreateEverythingExactly("templates/", "domain", fileRenamer, domainObj, utils.AppTemplates)
		if err != nil {
			return err
		}

	}

	return nil

}

func insertNewDomainName(filePath, domainName string) error {

	f, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {

		}
	}(f)

	isEmptyFile := true

	fileContent := ""
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		if line == "" {
			continue
		}

		isEmptyFile = false

		x := line
		if strings.HasPrefix(line, "-") {
			x = line[1:]
		}

		if x == domainName {
			return fmt.Errorf("domain name already exist")
		}

		fileContent += line
		fileContent += "\n"
	}
	if err := scanner.Err(); err != nil {
		return err
	}

	if isEmptyFile {
		fileContent += fmt.Sprintf("-%s", domainName)
	} else {
		fileContent += domainName
	}

	fileContent += "\n"

	return os.WriteFile(filePath, []byte(fileContent), 0644)
}
package gendomain

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/mirzaakhena/gogen/utils"
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

	packagePath := utils.GetPackagePath()

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
		PackagePath:    packagePath,
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

	err = utils.CreateGogenConfig(err, domainName)
	if err != nil {
		return err
	}

	// handle .gitignore
	{

		gitignoreFile := fmt.Sprintf("templates/domain/%s", "~gitignore")

		fileBytes, err := utils.AppTemplates.ReadFile(gitignoreFile)
		if err != nil {
			return err
		}

		_, err = utils.WriteFileIfNotExist(string(fileBytes), "./.gitignore", struct{}{})
		if err != nil {
			return err
		}

	}

	return nil

}

//func insertNewDomainName(filePath, domainName string) error {
//
//	f, err := os.Open(filePath)
//	if err != nil {
//		return err
//	}
//	defer func(f *os.File) {
//		err := f.Close()
//		if err != nil {
//
//		}
//	}(f)
//
//	isEmptyFile := true
//
//	fileContent := ""
//	scanner := bufio.NewScanner(f)
//	for scanner.Scan() {
//
//		line := strings.TrimSpace(scanner.Text())
//
//		if line == "" {
//			continue
//		}
//
//		isEmptyFile = false
//
//		x := line
//		if strings.HasPrefix(line, "-") {
//			x = line[1:]
//		}
//
//		if x == domainName {
//			return fmt.Errorf("domain name already exist")
//		}
//
//		fileContent += line
//		fileContent += "\n"
//	}
//	if err := scanner.Err(); err != nil {
//		return err
//	}
//
//	if isEmptyFile {
//		fileContent += fmt.Sprintf("-%s", domainName)
//	} else {
//		fileContent += domainName
//	}
//
//	fileContent += "\n"
//
//	return os.WriteFile(filePath, []byte(fileContent), 0644)
//}
//
//func Rewrite(srcDir, destDir string) {
//
//	// Walk through the source directory and copy each file and subdirectory to the destination directory
//	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
//		if err != nil {
//			fmt.Println(err)
//			return err
//		}
//
//		// Construct the destination path by replacing the source directory with the destination directory
//		destPath := filepath.Join(destDir, path[len(srcDir):])
//
//		// If the current path is a directory, create it in the destination directory
//		if info.IsDir() {
//			os.MkdirAll(destPath, info.Mode())
//			return nil
//		}
//
//		// If the current path is a file, copy it to the destination directory
//		srcFile, err := os.Open(path)
//		if err != nil {
//			fmt.Println(err)
//			return err
//		}
//		defer srcFile.Close()
//
//		destFile, err := os.Create(destPath)
//		if err != nil {
//			fmt.Println(err)
//			return err
//		}
//		defer destFile.Close()
//
//		_, err = io.Copy(destFile, srcFile)
//		if err != nil {
//			fmt.Println(err)
//			return err
//		}
//
//		return nil
//	})
//	if err != nil {
//		return
//	}
//}

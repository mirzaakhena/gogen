package gendomain

import (
	"bufio"
	"embed"
	"fmt"
	"io"
	"math/rand"
	"os"
	"path/filepath"
	"strings"
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

	_, err = utils.CreateFolderIfNotExist(".gogen")
	if err != nil {
		return err
	}

	gogenDomainFile := "./.gogen/domain"
	exist, err := utils.WriteFileIfNotExist(defaultDomain, gogenDomainFile, struct{}{})
	if err != nil {
		return err
	}

	//_, err = utils.CreateFolderIfNotExist(".gogen/templates/controller")
	//if err != nil {
	//	return err
	//}
	//
	//_, err = utils.CreateFolderIfNotExist(".gogen/templates/gateway")
	//if err != nil {
	//	return err
	//}

	{

		ff := FileAndFolders{
			Folders: map[string]int{},
			Files:   make([]string, 0),
		}
		err = readFolders(utils.AppTemplates, "templates/controllers", "templates/controllers", &ff)
		if err != nil {
			return err
		}

		for folder := range ff.Folders {

			folderName := fmt.Sprintf("%v/%v", ".gogen/templates/controllers", folder)

			err := os.MkdirAll(folderName, 0755)
			if err != nil {
				return err
			}

		}

		for _, fileRaw := range ff.Files {

			//fmt.Printf("file   >>>> %v/%v\n", ".gogen/templates/controllers", fileRaw)

			filename := fmt.Sprintf("%v/%v", ".gogen/templates/controllers", fileRaw)
			file, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("000 %v", err)
			}
			defer file.Close()

			fileInBytes, err := utils.AppTemplates.ReadFile("templates/controllers/" + fileRaw)
			if err != nil {
				return fmt.Errorf("111 %v", err)
			}

			_, err = file.WriteString(string(fileInBytes))
			if err != nil {
				return fmt.Errorf("222 %v", err)
			}

			// Sync the file to disk
			err = file.Sync()
			if err != nil {
				return fmt.Errorf("333 %v", err)
			}

		}

	}

	{

		ff := FileAndFolders{
			Folders: map[string]int{},
			Files:   make([]string, 0),
		}
		err = readFolders(utils.AppTemplates, "templates/gateway", "templates/gateway", &ff)
		if err != nil {
			return err
		}

		for folder := range ff.Folders {

			folderName := fmt.Sprintf("%v/%v", ".gogen/templates/gateway", folder)

			err := os.MkdirAll(folderName, 0755)
			if err != nil {
				return err
			}

		}

		for _, fileRaw := range ff.Files {

			//fmt.Printf("file   >>>> %v/%v\n", ".gogen/templates/gateway", fileRaw)

			filename := fmt.Sprintf("%v/%v", ".gogen/templates/gateway", fileRaw)
			file, err := os.Create(filename)
			if err != nil {
				return fmt.Errorf("000 %v", err)
			}
			defer file.Close()

			fileInBytes, err := utils.AppTemplates.ReadFile("templates/gateway/" + fileRaw)
			if err != nil {
				return fmt.Errorf("111 %v", err)
			}

			_, err = file.WriteString(string(fileInBytes))
			if err != nil {
				return fmt.Errorf("222 %v", err)
			}

			// Sync the file to disk
			err = file.Sync()
			if err != nil {
				return fmt.Errorf("333 %v", err)
			}

		}

	}

	if exist {
		_ = insertNewDomainName(gogenDomainFile, domainName)
	}

	gitignoreContent := `
.idea/
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

	//inFile, err := os.Open(".gogen/domain")
	//if err != nil {
	//	return err
	//}
	//defer func(inFile *os.File) {
	//	err := inFile.Close()
	//	if err != nil {
	//
	//	}
	//}(inFile)
	//
	//scanner := bufio.NewScanner(inFile)
	//for scanner.Scan() {
	//	domainNameInGogenFile := strings.TrimSpace(scanner.Text())
	//	if domainNameInGogenFile == "" {
	//		continue
	//	}
	//	if strings.HasPrefix(domainNameInGogenFile, "-") {
	//		domainNameInGogenFile = strings.ReplaceAll(domainNameInGogenFile, "-", "")
	//	}
	//	domainNameInGogenFile = strings.ToLower(domainNameInGogenFile)
	//	//_, err := utils.CreateFolderIfNotExist(fmt.Sprintf("domain_%s", domainNameInGogenFile))
	//	//if err != nil {
	//	//	return err
	//	//}
	//
	//	fileRenamer := map[string]string{
	//		"domainname": utils.LowerCase(domainNameInGogenFile),
	//	}
	//
	//	domainObj := ObjTemplate{DomainName: domainNameInGogenFile}
	//
	//	err = utils.CreateEverythingExactly("templates/", "domain", fileRenamer, domainObj, utils.AppTemplates)
	//	if err != nil {
	//		return err
	//	}
	//
	//}

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

func Rewrite(srcDir, destDir string) {

	// Walk through the source directory and copy each file and subdirectory to the destination directory
	err := filepath.Walk(srcDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println(err)
			return err
		}

		// Construct the destination path by replacing the source directory with the destination directory
		destPath := filepath.Join(destDir, path[len(srcDir):])

		// If the current path is a directory, create it in the destination directory
		if info.IsDir() {
			os.MkdirAll(destPath, info.Mode())
			return nil
		}

		// If the current path is a file, copy it to the destination directory
		srcFile, err := os.Open(path)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer srcFile.Close()

		destFile, err := os.Create(destPath)
		if err != nil {
			fmt.Println(err)
			return err
		}
		defer destFile.Close()

		_, err = io.Copy(destFile, srcFile)
		if err != nil {
			fmt.Println(err)
			return err
		}

		return nil
	})
	if err != nil {
		return
	}
}

type FileAndFolders struct {
	Folders map[string]int
	Files   []string
}

func readFolders(efs embed.FS, skip, path string, ff *FileAndFolders) error {

	dirs, err := efs.ReadDir(path)
	if err != nil {
		return err
	}

	for _, dir := range dirs {

		name := dir.Name()

		if dir.IsDir() {

			s := fmt.Sprintf("%s/%s", path, name)
			//fmt.Printf("found folder %s\n", s)

			for k := range ff.Folders {
				//fmt.Printf("k=%v\n", k)
				if strings.Contains(s, k) {
					//fmt.Printf("remove %v from %v\n", k, ff.Folders)
					delete(ff.Folders, k)
				}
			}

			ff.Folders[s[len(skip)+1:]] = 1

			err = readFolders(efs, skip, fmt.Sprintf("%s/%s", path, name), ff)
			if err != nil {
				return err
			}

		} else {
			s := fmt.Sprintf("%s/%s", path, name)
			//fmt.Printf("found file   %s\n", s)
			fileName := s[len(skip)+1:]
			//fmt.Printf("%v\n", fileName)
			ff.Files = append(ff.Files, fileName)
		}

	}

	return nil

}

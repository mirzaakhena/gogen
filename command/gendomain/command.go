package gendomain

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/mirzaakhena/gogen/utils/model"
	"math/rand"
	"os"
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

	err = CopyPasteFolder(".gogen/templates", "controller")
	if err != nil {
		return err
	}

	err = CopyPasteFolder(".gogen/templates", "gateway")
	if err != nil {
		return err
	}

	// handle gogen/config.json
	{
		gogenDomainFile := "./.gogen/gogenrc.json"

		data := model.GogenConfig{
			Domain:     domainName,
			Controller: "default",
			Gateway:    "default",
		}

		jsonInBytes, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			return err
		}
		_, err = utils.WriteFileIfNotExist(string(jsonInBytes), gogenDomainFile, struct{}{})
		if err != nil {
			return err
		}
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

func CopyPasteFolder(destFolder, sourceFolder string) error {

	ff := FileAndFolders{
		Folders: map[string]int{},
		Files:   make([]string, 0),
	}

	templateController := fmt.Sprintf("templates/%s", sourceFolder)
	err := readFolders(utils.AppTemplates, templateController, templateController, &ff)
	if err != nil {
		return err
	}

	for folder := range ff.Folders {
		err := os.MkdirAll(fmt.Sprintf("%s/%s/%v", destFolder, sourceFolder, folder), 0755)
		if err != nil {
			return err
		}
	}

	for _, fileRaw := range ff.Files {

		file, err := os.Create(fmt.Sprintf("%s/%s/%s", destFolder, sourceFolder, fileRaw))
		if err != nil {
			return err
		}

		fileInBytes, err := utils.AppTemplates.ReadFile(fmt.Sprintf("templates/%s/%s", sourceFolder, fileRaw))
		if err != nil {
			return err
		}

		_, err = file.WriteString(string(fileInBytes))
		if err != nil {
			return err
		}

		err = file.Sync()
		if err != nil {
			return err
		}

		err = file.Close()
		if err != nil {
			return err
		}
	}

	return nil
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

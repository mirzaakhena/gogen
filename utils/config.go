package utils

import (
	"embed"
	"encoding/json"
	"fmt"
	"github.com/mirzaakhena/gogen/utils/model"
	"os"
	"strings"
)

func CreateGogenConfig(err error, domainName string) error {
	err = CopyPasteFolder(".gogen/templates", "controller")
	if err != nil {
		return err
	}

	err = CopyPasteFolder(".gogen/templates", "gateway")
	if err != nil {
		return err
	}

	err = CopyPasteFolder(".gogen/templates", "crud")
	if err != nil {
		return err
	}

	// handle gogen/config.json
	{
		gogenDomainFile := "./.gogen/gogenrc.json"

		data := model.GogenConfig{
			Domain:     domainName,
			Controller: "gin",
			Gateway:    "simple",
			Crud:       "gin",
		}

		jsonInBytes, err := json.MarshalIndent(data, "", " ")
		if err != nil {
			return err
		}
		_, err = WriteFileIfNotExist(string(jsonInBytes), gogenDomainFile, struct{}{})
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
	err := readFolders(AppTemplates, templateController, templateController, &ff)
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

		fileInBytes, err := AppTemplates.ReadFile(fmt.Sprintf("templates/%s/%s", sourceFolder, fileRaw))
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

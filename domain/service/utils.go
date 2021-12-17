package service

import (
	"embed"
	"fmt"
	"github.com/mirzaakhena/gogen/gateway/prod"
	"github.com/mirzaakhena/gogen/infrastructure/templates"
	"github.com/mirzaakhena/gogen/infrastructure/util"
	"os"
	"strings"
	"text/template"
)

func CreateEverythingExactly(skip, path string, fileRenamer map[string]string, data interface{}) error {

	path = fmt.Sprintf("%s%s", skip, path)

	lenSkip := len(skip)

	ff := fileAndFolders{
		Folders: map[string]int{},
		Files:   make([]string, 0),
	}

	// TODO need to refactor later
	efs := templates.AppTemplates
	if strings.HasPrefix(skip, "crud") {
		efs = templates.CrudTemplates
	}

	if strings.HasPrefix(skip, "webapp") {
		efs = templates.WebappTemplates
	}

	if strings.HasPrefix(skip, "endtoend") {
		efs = templates.E2ETemplates
	}

	err := readFolders(efs, path, &ff)
	if err != nil {
		return err
	}

	i := strings.Index(path, skip)

	err = os.MkdirAll(path[i+lenSkip:], 0755)
	if err != nil {
		return err
	}

	for folder := range ff.Folders {

		s := replaceVariable(folder, fileRenamer)

		k := strings.Index(s, skip)

		err := os.MkdirAll(s[k+lenSkip:], 0755)
		if err != nil {
			return err
		}
	}

	for _, fileRaw := range ff.Files {

		file := replaceVariable(fileRaw, fileRenamer)

		i := strings.LastIndex(file, "/")
		nameFileWithExtOnly := fmt.Sprintf("%s", file[i+1:])

		if strings.HasPrefix(nameFileWithExtOnly, "~") {
			continue
		}

		j := strings.LastIndex(nameFileWithExtOnly, "._")
		nameFileWithoutUnderscore := fmt.Sprintf("%s/%s%s", file[:i], nameFileWithExtOnly[:j+1], nameFileWithExtOnly[j+2:])

		// skip the first path
		k := strings.Index(nameFileWithoutUnderscore, skip)

		if util.IsFileExist(nameFileWithoutUnderscore[k+lenSkip:]) {
			continue
		}

		fileOut, err := os.Create(nameFileWithoutUnderscore[k+lenSkip:])
		if err != nil {
			return err
		}

		templateData, err := efs.ReadFile(fileRaw)
		if err != nil {
			return err
		}

		tpl, err := template.
			New("something").
			Funcs(prod.FuncMap).
			Parse(string(templateData))

		template.New("").Lookup("")

		if err != nil {
			return err
		}

		err = tpl.Execute(fileOut, data)
		if err != nil {
			return err
		}

	}

	//fmt.Printf("%v\n", ff)

	return nil
}

func replaceVariable(folder string, fileRenamer map[string]string) string {
	s := folder
	for k, v := range fileRenamer {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%v}", k), v)
	}
	return s
}

type fileAndFolders struct {
	Folders map[string]int
	Files   []string
}

func readFolders(efs embed.FS, path string, ff *fileAndFolders) error {

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

			ff.Folders[s] = 1

			err = readFolders(efs, fmt.Sprintf("%s/%s", path, name), ff)
			if err != nil {
				return err
			}

		} else {
			s := fmt.Sprintf("%s/%s", path, name)
			//fmt.Printf("found file   %s\n", s)
			ff.Files = append(ff.Files, s)
		}

	}

	return nil

}

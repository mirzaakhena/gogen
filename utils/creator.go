package utils

import (
	"embed"
	"fmt"
	"os"
	"strings"
	"text/template"
)

func CreateEverythingExactly(rootToSkip, pathUnder string, fileRenamer map[string]string, data any, efs embed.FS) error {
	return createEverythingImpl{}.CreateEverythingExactly(rootToSkip+pathUnder, fileRenamer, data, efs)
}

type createEverythingImpl struct{}

func (x createEverythingImpl) CreateEverythingExactly(pathUnder string, fileRenamer map[string]string, data any, efs embed.FS) error {

	ff := FileAndFolders{
		Folders: map[string]int{},
		Files:   make([]string, 0),
	}

	err := x.readFolders(efs, pathUnder, pathUnder, &ff)
	if err != nil {
		return err
	}

	for folder := range ff.Folders {

		s := replaceVariable(folder, fileRenamer)

		err := os.MkdirAll(s, 0755)
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

		var nameFileWithoutUnderscore string

		if i == -1 {

			if file == "Dockerfile" {
				nameFileWithoutUnderscore = file
			} else {
				nameFileWithoutUnderscore = fmt.Sprintf("%s%s", nameFileWithExtOnly[:j+1], nameFileWithExtOnly[j+2:])
			}

		} else {

			nameFileWithoutUnderscore = fmt.Sprintf("%s/%s%s", file[:i], nameFileWithExtOnly[:j+1], nameFileWithExtOnly[j+2:])

		}

		if IsFileExist(nameFileWithoutUnderscore) {
			continue
		}

		fileOut, err := os.Create(nameFileWithoutUnderscore)
		if err != nil {
			return err
		}

		templateData, err := efs.ReadFile(pathUnder + "/" + fileRaw)
		if err != nil {
			return err
		}

		tpl, err := template.
			New("something").
			Funcs(FuncMap).
			Parse(string(templateData))

		template.New("").Lookup("")

		if err != nil {
			return err
		}

		if data == nil {
			data = struct {
			}{}
		}

		err = tpl.Execute(fileOut, data)
		if err != nil {
			return err
		}

		if strings.HasSuffix(nameFileWithoutUnderscore, ".go") {

			// temporary handling need to fixed later
			if strings.HasPrefix(nameFileWithExtOnly, "main") {
				continue
			}

			// reformat the file
			err = Reformat(nameFileWithoutUnderscore, nil)
			if err != nil {
				return err
			}

		}

	}

	//fmt.Printf("%v\n", ff)

	return nil
}

//func (x createEverythingImpl) replaceVariable(folder string, fileRenamer map[string]string) string {
//	if fileRenamer == nil {
//		return folder
//	}
//	s := folder
//	for k, v := range fileRenamer {
//		s = strings.ReplaceAll(s, fmt.Sprintf("${%v}", k), v)
//	}
//	return s
//}

type FileAndFolders struct {
	Folders map[string]int
	Files   []string
}

func (x createEverythingImpl) readFolders(efs embed.FS, skip, path string, ff *FileAndFolders) error {

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

			err = x.readFolders(efs, skip, fmt.Sprintf("%s/%s", path, name), ff)
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

func replaceVariable(folder string, fileRenamer map[string]string) string {
	if fileRenamer == nil {
		return folder
	}
	s := folder
	for k, v := range fileRenamer {
		s = strings.ReplaceAll(s, fmt.Sprintf("${%v}", k), v)
	}
	return s
}

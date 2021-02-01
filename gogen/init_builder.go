package gogen

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

type InitBuilderRequest struct {
	FolderPath string
	GomodPath  string
}

type initBuilder struct {
	InitBuilderRequest InitBuilderRequest
}

func NewInit(req InitBuilderRequest) Generator {
	return &initBuilder{InitBuilderRequest: req}
}

func (d *initBuilder) Generate() error {
	folderPath := d.InitBuilderRequest.FolderPath

	source := DefaultTemplatePath("")
	destination := fmt.Sprintf("%s/.gogen/templates/default", folderPath)
	if err := copyDir(source, destination); err != nil {
		return err
	}
	data, err := ioutil.ReadFile(fmt.Sprintf("%s/src/github.com/mirzaakhena/gogen/config.json", GetGopath()))
	if err != nil {
		return err
	}

	_ = ioutil.WriteFile(fmt.Sprintf("%s/.gogen/config.json", folderPath), data, 0644)

	return nil
}

func copyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = copyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = copyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func copyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// copyDir and copyFile is based on
// https://gist.github.com/r0l1/92462b38df26839a3ca324697c8cba04

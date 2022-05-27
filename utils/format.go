package utils

import (
	"golang.org/x/tools/imports"
	"io/ioutil"
)

func Reformat(goFilename string, bytes []byte) error {

	// reformat the import
	newBytes, err := imports.Process(goFilename, bytes, nil)
	if err != nil {
		return err
	}

	// rewrite it
	if err := ioutil.WriteFile(goFilename, newBytes, 0644); err != nil {
		return err
	}

	return nil
}

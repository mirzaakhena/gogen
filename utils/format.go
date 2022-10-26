package utils

import (
	"golang.org/x/tools/imports"
	"os"
)

func Reformat(goFilename string, bytes []byte) error {

	// reformat the import
	newBytes, err := imports.Process(goFilename, bytes, nil)
	if err != nil {
		return err
	}

	// rewrite it
	if err := os.WriteFile(goFilename, newBytes, 0644); err != nil {
		return err
	}

	return nil
}

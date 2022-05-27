package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

// GetPackageName ...
func GetPackageName(rootFolderName string) string {
	i := strings.LastIndex(rootFolderName, "/")
	return rootFolderName[i+1:]
}

func GetDefaultDomain() string {

	var defaultDomain string

	file, err := os.Open("./.gogen/domain")
	if err != nil {
		fmt.Printf(".gogen/domain is not found. Please run 'gogen init' first\n")
		os.Exit(1)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	defer func() {
		if err := scanner.Err(); err != nil {
			panic(err.Error())
		}
	}()
	found := false
	for scanner.Scan() {
		row := scanner.Text()
		if strings.HasPrefix(row, "-") {

			if found {
				fmt.Printf("Found multiple selected domain. Put just one '-' in front of domain name. \n")
				os.Exit(1)
			}

			i := strings.Index(row, "-")
			defaultDomain = strings.TrimSpace(row[i+1:])
			found = true
		}
	}

	if !found {
		fmt.Printf("No domain selected. Please select one of domain by put '-' in front of domain name\n")
		os.Exit(1)
	}

	return strings.ToLower(defaultDomain)

}

func GetPackagePath() string {

	var gomodPath string

	file, err := os.Open("go.mod")
	if err != nil {
		fmt.Printf("go.mod is not found. Please create it with command `go mod init your/path/project`\n")
		os.Exit(1)
	}
	defer func() {
		err = file.Close()
		if err != nil {
			return
		}
	}()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		row := scanner.Text()
		if strings.HasPrefix(row, "module") {
			moduleRow := strings.Split(row, " ")
			if len(moduleRow) > 1 {
				gomodPath = moduleRow[1]
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err.Error())
	}

	return strings.Trim(gomodPath, "\"")

}

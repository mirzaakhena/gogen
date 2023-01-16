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

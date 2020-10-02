package gogen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type test struct {
}

func NewTest() Generator {
	return &test{}
}

func (d *test) Generate(args ...string) error {

	if len(args) < 3 {
		return fmt.Errorf("please define test usecase. ex: `gogen test CreateOrder`")
	}

	usecaseName := args[2]

	ds := Test{}
	ds.UsecaseName = usecaseName
	ds.PackagePath = GetPackagePath()

	file, err := os.Open(fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(usecaseName)))
	if err != nil {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	state := 0
	for scanner.Scan() {
		if strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sOutport interface {", usecaseName)) {
			state = 1
		} else //
		if state == 1 {
			if strings.HasPrefix(scanner.Text(), "}") {
				state = 2
				break
			} else {
				completeMethod := strings.TrimSpace(scanner.Text())
				methodNameOnly := strings.Split(completeMethod, "(")[0]
				ds.OutportMethods = append(ds.OutportMethods, methodNameOnly)
			}
		}
	}

	if state == 0 {
		return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
	}

	_ = WriteFileIfNotExist(
		"usecase/usecaseName/interactor_test._go",
		fmt.Sprintf("usecase/%s/interactor_test.go", strings.ToLower(usecaseName)),
		ds,
	)

	GenerateMock(ds.PackagePath, ds.UsecaseName)

	return nil
}

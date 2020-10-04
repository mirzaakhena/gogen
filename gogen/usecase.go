package gogen

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type usecase struct {
}

func NewUsecase() Generator {
	return &usecase{}
}

func (d *usecase) Generate(args ...string) error {

	if len(args) < 4 {
		return fmt.Errorf("please define usecase name and type (command/query). ex: `gogen usecase command CreateOrder`")
	}

	usecaseType := args[2]

	usecaseName := args[3]

	packagePath := GetPackagePath()

	uc := Usecase{
		Name:        usecaseName,
		PackagePath: packagePath,
	}

	CreateFolder("usecase/%s/port", strings.ToLower(uc.Name))

	if usecaseType == "command" {

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/inport-command._go",
			fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(uc.Name)),
			uc,
		)

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/outport-command._go",
			fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(uc.Name)),
			uc,
		)

		{
			file, err := os.Open(fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", usecaseName)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {
						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						uc.InportRequestFields = append(uc.InportRequestFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
							Type: strings.TrimSpace(fieldWithType[1]),
						})
					}
				}
			}
			if state == 0 {
				return fmt.Errorf("not found Request struct")
			}
		}

		{
			file, err := os.Open(fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", usecaseName)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {
						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						uc.InportResponseFields = append(uc.InportResponseFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
							Type: strings.TrimSpace(fieldWithType[1]),
						})
					}
				}
			}
			if state == 0 {
				return fmt.Errorf("not found Response struct")
			}
		}

		{
			file, err := os.Open(fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sOutport interface {", usecaseName)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {
						completeMethod := strings.TrimSpace(scanner.Text())
						methodNameOnly := strings.Split(completeMethod, "(")[0]
						uc.Outports = append(uc.Outports, &Outport{
							Name: methodNameOnly,
						})
					}
				}
			}

			if state == 0 {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
		}

		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", ot.Name)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {

						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						ot.RequestFields = append(ot.RequestFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
						})

					}
				}
			}

			if state == 0 {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
		}

		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", ot.Name)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {

						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						ot.ResponseFields = append(ot.ResponseFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
						})

					}
				}
			}

			if state == 0 {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
		}

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor-command._go",
			fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(uc.Name)),
			uc,
		)

	} else //

	if usecaseType == "query" {

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/inport-query._go",
			fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(uc.Name)),
			uc,
		)

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/outport-query._go",
			fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(uc.Name)),
			uc,
		)

		{
			file, err := os.Open(fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", usecaseName)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {
						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						uc.InportRequestFields = append(uc.InportRequestFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
							Type: strings.TrimSpace(fieldWithType[1]),
						})
					}
				}
			}
			if state == 0 {
				return fmt.Errorf("not found Request struct")
			}
		}

		{
			file, err := os.Open(fmt.Sprintf("usecase/%s/port/inport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", usecaseName)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {
						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						uc.InportResponseFields = append(uc.InportResponseFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
							Type: strings.TrimSpace(fieldWithType[1]),
						})
					}
				}
			}
			if state == 0 {
				return fmt.Errorf("not found Response struct")
			}
		}

		{
			file, err := os.Open(fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sOutport interface {", usecaseName)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {
						completeMethod := strings.TrimSpace(scanner.Text())
						methodNameOnly := strings.Split(completeMethod, "(")[0]
						uc.Outports = append(uc.Outports, &Outport{
							Name: methodNameOnly,
						})
					}
				}
			}

			if state == 0 {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
		}

		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", ot.Name)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {

						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						ot.RequestFields = append(ot.RequestFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
						})

					}
				}
			}

			if state == 0 {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
		}

		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("usecase/%s/port/outport.go", strings.ToLower(usecaseName)))
			if err != nil {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", ot.Name)) {
					state = 1
				} else //
				if state == 1 {
					if strings.HasPrefix(scanner.Text(), "}") {
						state = 2
						break
					} else {

						completeFieldWithType := strings.TrimSpace(scanner.Text())
						if len(completeFieldWithType) == 0 {
							continue
						}
						fieldWithType := strings.SplitN(completeFieldWithType, " ", 2)
						ot.ResponseFields = append(ot.ResponseFields, &NameType{
							Name: strings.TrimSpace(fieldWithType[0]),
						})

					}
				}
			}

			if state == 0 {
				return fmt.Errorf("not found usecase %s. You need to create it first by call 'gogen usecase %s' ", usecaseName, usecaseName)
			}
		}

		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor-query._go",
			fmt.Sprintf("usecase/%s/interactor.go", strings.ToLower(uc.Name)),
			uc,
		)

	} else //

	{
		return fmt.Errorf("use type `command` or `query`")
	}

	return nil
}

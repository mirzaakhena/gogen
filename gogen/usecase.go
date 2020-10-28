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

	return GenerateUsecase(UsecaseRequest{
		UsecaseType: args[2],
		UsecaseName: args[3],
		FolderPath:  ".",
		// OutportMethods: []string{"DoSomething"},
	})
}

type UsecaseRequest struct {
	UsecaseType    string
	UsecaseName    string
	OutportMethods []string
	FolderPath     string
}

func GenerateUsecase(req UsecaseRequest) error {

	packagePath := GetPackagePath()

	var folderImport string
	if req.FolderPath != "." {
		folderImport = fmt.Sprintf("/%s", req.FolderPath)
	}

	uc := Usecase{
		Name:        req.UsecaseName,
		PackagePath: packagePath,
		Directory:   folderImport,
	}

	firstTime := false

	if !IsExist(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(uc.Name))) {
		if len(req.OutportMethods) == 0 {
			req.OutportMethods = []string{"DoSomething"}
			firstTime = true
		}
	}

	// set outport methods
	{
		outports := []*Outport{}
		for _, methodName := range req.OutportMethods {
			outports = append(outports, &Outport{
				Name:           methodName,
				RequestFields:  nil,
				ResponseFields: nil,
			})
		}
		uc.Outports = outports
	}

	// Create Port Folder
	CreateFolder("%s/usecase/%s/port", req.FolderPath, strings.ToLower(uc.Name))

	if req.UsecaseType == "command" {

		// Create inport file
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/inport-command._go",
			fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(uc.Name)),
			uc,
		)

		// Create outport file
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/outport._go",
			fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(uc.Name)),
			uc,
		)

		// Read inport request field
		{
			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error1. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", req.UsecaseName)) {
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

		// Read inport response field
		{
			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error2. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", req.UsecaseName)) {
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

		// read outport methods name
		if !firstTime {
			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error3. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sOutport interface {", req.UsecaseName)) {
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
				return fmt.Errorf("error4. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
		}

		// read outport request field
		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error5. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
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
				return fmt.Errorf("error6. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
		}

		// read outport response field
		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error7. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
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
				return fmt.Errorf("error8. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
		}

		// create interactor file
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor-command._go",
			fmt.Sprintf("%s/usecase/%s/interactor.go", req.FolderPath, strings.ToLower(uc.Name)),
			uc,
		)

	} else //

	if req.UsecaseType == "query" {

		// create inport folder
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/inport-query._go",
			fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(uc.Name)),
			uc,
		)

		// create outport folder
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/port/outport._go",
			fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(uc.Name)),
			uc,
		)

		// read inport request field
		{
			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error9. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sRequest struct {", req.UsecaseName)) {
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

		// read inport response field
		{
			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/inport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error10. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sResponse struct {", req.UsecaseName)) {
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

		// read outport methods name
		if !firstTime {
			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error11. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
			defer file.Close()

			scanner := bufio.NewScanner(file)
			scanner.Split(bufio.ScanLines)

			state := 0
			for scanner.Scan() {
				if state == 0 && strings.HasPrefix(scanner.Text(), fmt.Sprintf("type %sOutport interface {", req.UsecaseName)) {
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
				return fmt.Errorf("error12. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
		}

		// read outport request fields
		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error13. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
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
				return fmt.Errorf("error14. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
		}

		// read outport response fields
		for _, ot := range uc.Outports {

			file, err := os.Open(fmt.Sprintf("%s/usecase/%s/port/outport.go", req.FolderPath, strings.ToLower(req.UsecaseName)))
			if err != nil {
				return fmt.Errorf("error15. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
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
				return fmt.Errorf("error16. not found usecase %s. You need to create it first by call 'gogen usecase %s' ", req.UsecaseName, req.UsecaseName)
			}
		}

		// create interactor file
		_ = WriteFileIfNotExist(
			"usecase/usecaseName/interactor-query._go",
			fmt.Sprintf("%s/usecase/%s/interactor.go", req.FolderPath, strings.ToLower(uc.Name)),
			uc,
		)

	} else //

	{
		return fmt.Errorf("use type `command` or `query`")
	}

	return nil
}

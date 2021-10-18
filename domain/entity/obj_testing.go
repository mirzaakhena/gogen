package entity

import (
	"fmt"
  "github.com/mirzaakhena/gogen/domain/vo"
)

// ObjTesting depend on (which) usecase that want to be tested
type ObjTesting struct {
	TestName   vo.Naming
	ObjUsecase ObjUsecase
}

// ObjDataTesting ...
type ObjDataTesting struct {
	PackagePath string
	UsecaseName string
	TestName    string
	Methods     vo.OutportMethods
}

// NewObjTesting ...
func NewObjTesting(testName string, objUsecase ObjUsecase) (*ObjTesting, error) {

	var obj ObjTesting
	obj.TestName = vo.Naming(testName)
	obj.ObjUsecase = objUsecase

	return &obj, nil
}

// GetData ...
func (o ObjTesting) GetData(PackagePath string, outportMethods vo.OutportMethods) *ObjDataTesting {
	return &ObjDataTesting{
		PackagePath: PackagePath,
		UsecaseName: o.ObjUsecase.UsecaseName.String(),
		TestName:    o.TestName.LowerCase(),
		Methods:     outportMethods,
	}
}

// GetTestFileName ...
func (o ObjTesting) GetTestFileName() string {
	return fmt.Sprintf("%s/testcase_%s_test.go", o.ObjUsecase.GetUsecaseRootFolderName(), o.TestName.LowerCase())
}

func (o ObjTesting) InjectToTest(injectedCode string) ([]byte, error) {
	return InjectCodeAtTheEndOfFile(o.GetTestFileName(), injectedCode)
}
package entity

import (
	"fmt"
  "github.com/mirzaakhena/gogen/application/apperror"
  "github.com/mirzaakhena/gogen/domain/vo"
)

// ObjError depend on (which) usecase that want to be tested
type ObjError struct {
	ErrorName vo.Naming
}

// ObjDataError is object that used in template
type ObjDataError struct {
	ErrorName string
}

// NewObjError Constructor
func NewObjError(errorName string) (*ObjError, error) {

	if errorName == "" {
		return nil, apperror.ErrorNameMustNotEmpty
	}

	var obj ObjError
	obj.ErrorName = vo.Naming(errorName)

	return &obj, nil
}

// TODO add function is Exist

// GetData ...
func (o ObjError) GetData() *ObjDataError {
	return &ObjDataError{
		ErrorName: o.ErrorName.String(),
	}
}

// GetErrorRootFolderName ...
func (o ObjError) GetErrorRootFolderName() string {
	return fmt.Sprintf("application/apperror")
}

// GetErrorEnumFileName ...
func (o ObjError) GetErrorEnumFileName() string {
	return fmt.Sprintf("%s/error_enum.go", o.GetErrorRootFolderName())
}

// GetErrorFuncFileName ...
func (o ObjError) GetErrorFuncFileName() string {
	return fmt.Sprintf("%s/error_func.go", o.GetErrorRootFolderName())
}

// InjectCode ...
func (o ObjError) InjectCode(templateCode string) ([]byte, error) {
	return InjectCodeAtTheEndOfFile(o.GetErrorEnumFileName(), templateCode)
}

package entity

import (
  "github.com/mirzaakhena/gogen/domain/vo"
)

// ObjValueObject ...
type ObjValueObject struct {
  ValueObjectName vo.Naming
  FieldNames      []string
}

// ObjDataValueObject ...
type ObjDataValueObject struct {
  ValueObjectName string
  FieldNames      []string
}

// NewObjValueObject ...
func NewObjValueObject(entityName string, fieldNames []string) (*ObjValueObject, error) {

  var obj ObjValueObject
  obj.ValueObjectName = vo.Naming(entityName)
  obj.FieldNames = fieldNames

  return &obj, nil
}

// GetData ...
func (o ObjValueObject) GetData() *ObjDataValueObject {
  return &ObjDataValueObject{
    ValueObjectName: o.ValueObjectName.String(),
    FieldNames:      o.FieldNames,
  }
}

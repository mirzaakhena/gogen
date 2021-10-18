package entity

import (
  "github.com/mirzaakhena/gogen/domain/vo"
)

// ObjValueString ...
type ObjValueString struct {
  ValueStringName vo.Naming
}

// ObjDataValueString ...
type ObjDataValueString struct {
  ValueStringName string
}

// NewObjValueString ...
func NewObjValueString(entityName string) (*ObjValueString, error) {

  var obj ObjValueString
  obj.ValueStringName = vo.Naming(entityName)

  return &obj, nil
}

// GetData ...
func (o ObjValueString) GetData() *ObjDataValueString {
  return &ObjDataValueString{
    ValueStringName: o.ValueStringName.String(),
  }
}

package entity

// ObjCrud ...
type ObjCrud struct {
	Entity *ObjEntity
}

// ObjDataCrud ...
type ObjDataCrud struct {
	EntityName  string
	PackagePath string
}

// NewObjCrud ...
func NewObjCrud(entityName string) (*ObjCrud, error) {

	var obj ObjCrud

	ent, err := NewObjEntity(entityName)
	if err != nil {
		return nil, err
	}

	obj.Entity = ent

	return &obj, nil
}

// GetData ...
func (o ObjCrud) GetData(PackagePath string) *ObjDataCrud {
	return &ObjDataCrud{
		EntityName:  o.Entity.EntityName.String(),
		PackagePath: PackagePath,
	}
}

package entity

// ObjEndToEnd ...
type ObjEndToEnd struct {
	Usecase *ObjUsecase
	Entity  *ObjEntity
}

// ObjDataEndToEnd ...
type ObjDataEndToEnd struct {
	UsecaseName string
	EntityName  string
	PackagePath string
}

// NewObjEndToEnd ...
func NewObjEndToEnd(usecaseName, entityName string) (*ObjEndToEnd, error) {

	var obj ObjEndToEnd

	us, err := NewObjUsecase(usecaseName)
	if err != nil {
		return nil, err
	}

	obj.Usecase = us

	ent, err := NewObjEntity(entityName)
	if err != nil {
		return nil, err
	}

	obj.Entity = ent

	return &obj, nil
}

// GetData ...
func (o ObjEndToEnd) GetData(PackagePath string) *ObjDataEndToEnd {
	return &ObjDataEndToEnd{
		UsecaseName: o.Usecase.UsecaseName.String(),
		EntityName:  o.Entity.EntityName.String(),
		PackagePath: PackagePath,
	}
}

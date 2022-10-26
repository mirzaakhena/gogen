package database

import "reflect"

type GetAllParam struct {
	Page   int64
	Size   int64
	Sort   map[string]any
	Filter map[string]any
}

func (g GetAllParam) SetPage(page int64) GetAllParam {
	g.Page = page
	return g
}

func (g GetAllParam) SetSize(size int64) GetAllParam {
	g.Size = size
	return g
}

func (g GetAllParam) SetSort(field string, sort any) GetAllParam {
	g.Sort[field] = sort
	return g
}

func (g GetAllParam) SetFilter(key string, value any) GetAllParam {

	if reflect.ValueOf(value).String() == "" {
		return g
	}

	g.Filter[key] = value
	return g
}

func NewDefaultParam() GetAllParam {
	return GetAllParam{
		Page:   1,
		Size:   2000,
		Sort:   map[string]any{},
		Filter: map[string]any{},
	}
}

// =======================================

type InsertOrUpdateRepo[T any] interface {
	InsertOrUpdate(obj *T) error
}

type InsertManyRepo[T any] interface {
	InsertMany(objs ...*T) error
}

type GetOneRepo[T any] interface {
	GetOne(filter map[string]any, result *T) error
}

type GetAllRepo[T any] interface {
	GetAll(param GetAllParam, results *[]*T) (int64, error)
}

type GetAllEachItemRepo[T any] interface {
	GetAllEachItem(param GetAllParam, resultEachItem func(result T)) (int64, error)
}

type DeleteRepo[T any] interface {
	Delete(filter map[string]any) error
}

type Repository[T any] interface {
	InsertOrUpdateRepo[T]
	InsertManyRepo[T]
	GetOneRepo[T]
	GetAllRepo[T]
	GetAllEachItemRepo[T]
	DeleteRepo[T]
	//GetCollection() *mongo.Collection

	GetTypeName() string
}

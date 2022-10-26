package gogen

import (
	"reflect"
	"strings"
)

type iUsecase interface {
	AddUsecase(inports ...any) iUsecase
}

type RegisterRouterHandler[T any] interface {
	iUsecase
	RegisterRouter(h T)
}

type BaseController struct {
	inportObjs map[any]any
}

func NewBaseController() *BaseController {
	return &BaseController{
		inportObjs: map[any]any{},
	}
}

func (r *BaseController) GetUsecase(nameStructType any) any {

	x := reflect.TypeOf(nameStructType).String()

	//fmt.Printf("keluar>>>> %v\n", x[:strings.Index(x, ".")])

	packageName := x[:strings.Index(x, ".")]

	uc, ok := r.inportObjs[packageName]
	if !ok {
		return nil
	}
	return uc
}

func (r *BaseController) AddUsecase(inports ...any) iUsecase {
	for _, inport := range inports {
		x := reflect.ValueOf(inport).Elem().Type().String()
		packagePath := x[:strings.Index(x, ".")]
		r.inportObjs[packagePath] = inport
	}
	return r
}

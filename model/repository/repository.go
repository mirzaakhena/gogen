package repository

import (
	"context"
	"github.com/mirzaakhena/gogen/model/entity"
)

type FindAllObjUsecasesRepo interface {
	FindAllObjUsecases(ctx context.Context, objController *entity.ObjController) ([]*entity.ObjUsecase, error)
}

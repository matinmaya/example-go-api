package module

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/pkg/paginator"
	"reapp/pkg/queryfilter"
)

type ServiceAdapter[T Identifiable[TID], TID UintID] struct {
	svc IService[T, TID]
}

func NewServiceAdapter[T Identifiable[TID], TID UintID](svc IService[T, TID]) ServiceAdapter[T, TID] {
	return ServiceAdapter[T, TID]{svc: svc}
}

func (a ServiceAdapter[T, TID]) Create(db *gorm.DB, model *T) error {
	return a.svc.Create(db, model)
}

func (a ServiceAdapter[T, TID]) GetByID(db *gorm.DB, id uint64) (*T, error) {
	return a.svc.GetByID(db, TID(id))
}

func (a ServiceAdapter[T, TID]) Update(db *gorm.DB, model *T) error {
	return a.svc.Update(db, model)
}

func (a ServiceAdapter[T, TID]) Delete(db *gorm.DB, id uint64) error {
	return a.svc.Delete(db, TID(id))
}

func (a ServiceAdapter[T, TID]) List(ctx *gin.Context, db *gorm.DB, pg *paginator.Pagination[T], filterFields []queryfilter.FilterField) error {
	return a.svc.List(ctx, db, pg, filterFields)
}

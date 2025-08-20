package basehandler

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/pkg/paginator"
	"reapp/pkg/queryfilter"
)

type IServiceLister[T any] interface {
	List(ctx *gin.Context, db *gorm.DB, pagination *paginator.Pagination[T], filterFields []queryfilter.FilterField) error
}

type IServiceGetter[T any] interface {
	GetAll(db *gorm.DB) ([]T, error)
}

type IServiceGetterDetail[T any] interface {
	GetDetail(db *gorm.DB, id uint64) (*T, error)
}

type IServiceCreator[T any] interface {
	Create(db *gorm.DB, model *T) error
}

type IServiceUpdater[T any] interface {
	Update(db *gorm.DB, model *T) error
	GetByID(db *gorm.DB, id uint64) (*T, error)
}

type IServiceDeleter interface {
	Delete(db *gorm.DB, id uint64) error
}

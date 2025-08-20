package module

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/pkg/paginator"
	"reapp/pkg/queryfilter"
)

type UintID interface {
	~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64
}

type Identifiable[TID UintID] interface {
	GetID() TID
}

type IService[T Identifiable[TID], TID UintID] interface {
	Create(db *gorm.DB, model *T) error
	GetByID(db *gorm.DB, id TID) (*T, error)
	Update(db *gorm.DB, model *T) error
	Delete(db *gorm.DB, id TID) error
	List(ctx *gin.Context, db *gorm.DB, pg *paginator.Pagination[T], filterFields []queryfilter.FilterField) error
}

type IServiceAdapter[T Identifiable[TID], TID UintID] interface {
	Create(db *gorm.DB, model *T) error
	GetByID(db *gorm.DB, id uint64) (*T, error)
	Update(db *gorm.DB, model *T) error
	Delete(db *gorm.DB, id uint64) error
	List(ctx *gin.Context, db *gorm.DB, pagination *paginator.Pagination[T], filterFields []queryfilter.FilterField) error
}

type IHandler[T Identifiable[TID], TID UintID, TDTO any, TQ any] interface {
	Create(ctx *gin.Context)
	Update(ctx *gin.Context)
	Delete(ctx *gin.Context)
	List(ctx *gin.Context)
}

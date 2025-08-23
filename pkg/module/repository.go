package module

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"reapp/pkg/paginator"
	"reapp/pkg/queryfilter"
	"reapp/pkg/services/rediservice"
)

type Repository[T IWithID[TID], TID IUintID] struct {
	namespace string
}

func NewRepository[T IWithID[TID], TID IUintID](namespace string) *Repository[T, TID] {
	return &Repository[T, TID]{namespace: namespace}
}

func (r *Repository[T, TID]) Create(db *gorm.DB, model *T) error {
	go rediservice.ClearCacheOfRepository(r.namespace)
	return db.Create(model).Error
}

func (r *Repository[T, TID]) GetByID(db *gorm.DB, id TID) (*T, error) {
	var model T
	err := db.First(&model, id).Error
	return &model, err
}

func (r *Repository[T, TID]) Update(db *gorm.DB, model *T) error {
	go rediservice.ClearCacheOfRepository(r.namespace)
	return db.Save(model).Error
}

func (r *Repository[T, TID]) Delete(db *gorm.DB, id TID) error {
	go rediservice.ClearCacheOfRepository(r.namespace)
	var model T
	return db.Delete(&model, id).Error
}

func (r *Repository[T, TID]) List(
	ctx *gin.Context, db *gorm.DB,
	pg *paginator.Pagination[T],
	filterFields []queryfilter.FilterField,
) error {
	var models []T
	var model T
	scopes := paginator.Paginate(db, r.namespace, &model, pg, filterFields)

	collectionKey := "list"
	if err := rediservice.CacheOfRepository(r.namespace, collectionKey, pg.GetListCacheKey(), &models); err != nil {
		err := db.Scopes(scopes).Find(&models).Error
		if err != nil {
			return err
		}

		rediservice.SetCacheOfRepository(r.namespace, collectionKey, pg.GetListCacheKey(), models)
	}

	pg.SetRows(models)
	return nil
}

func (r *Repository[T, TID]) SetNamespace(name string) {
	if r.namespace != "" {
		r.namespace = name
	}
}

func (r *Repository[T, TID]) IsModelFetched(model T) bool {
	return model.IsCreated() && model.GetID() != 0
}

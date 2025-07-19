package service

import (
	"reapp/pkg/filterscopes"
	"reapp/pkg/paginator"

	"gorm.io/gorm"
)

type Lister interface {
	List(db *gorm.DB, pagination *paginator.Pagination, filters []filterscopes.QueryFilter) error
}

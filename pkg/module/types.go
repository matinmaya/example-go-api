package module

import "gorm.io/gorm"

type TBusinessLogic[T any] func(db *gorm.DB, model *T) error

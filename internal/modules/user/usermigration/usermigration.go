package usermigration

import (
	"reapp/internal/modules/user/permmodel"
	"reapp/internal/modules/user/rolemodel"
	"reapp/internal/modules/user/usermodel"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) {
	db.AutoMigrate(
		&permmodel.Permission{},
		&rolemodel.Role{},
		&rolemodel.RolePermission{},
		&usermodel.User{},
		&usermodel.UserRole{},
	)
}

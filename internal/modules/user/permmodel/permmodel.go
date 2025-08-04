package permmodel

import (
	"reapp/pkg/base/basemodel"
	"reapp/pkg/validators"
)

type Permission struct {
	basemodel.PrimaryKey
	Name        string `json:"name" gorm:"unique;not null;type:varchar(100)" validate:"required,max=100,unique=sys_permissions?id"`
	Description string `json:"description" gorm:"type:varchar(255)" validate:"max=255"`
	basemodel.SoftFields
	validators.ValidateUniqueScope
}

func (Permission) TableName() string {
	return "sys_permissions"
}

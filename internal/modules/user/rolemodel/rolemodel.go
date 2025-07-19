package rolemodel

import (
	"reapp/internal/modules/user/permmodel"
	"reapp/pkg/basemodel"
	"reapp/pkg/validators"
)

type Role struct {
	basemodel.SmallPrimaryKey
	Name        string                 `json:"name" gorm:"unique;not null;type:varchar(50);" validate:"required,max=50,unique=sys_roles?id"`
	Description string                 `json:"description" gorm:"type:varchar(255);" validate:"max=255"`
	Status      uint8                  `json:"status" gorm:"not null;default:0;"`
	Permissions []permmodel.Permission `json:"permissions,omitempty" gorm:"many2many:sys_role_permission;"`
	basemodel.SoftFields
	validators.ValidateScopeUnique
}

func (Role) TableName() string {
	return "sys_roles"
}

type RolePermission struct {
	RoleID       uint16                   `json:"role_id" gorm:"primaryKey;column:role_id;not null"`
	PermissionID uint32                   `json:"permission_id" gorm:"primaryKey;column:permission_id;not null"`
	Role         Role                     `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Permission   permmodel.Permission     `gorm:"foreignKey:PermissionID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt    basemodel.DateTimeFormat `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "sys_role_permission"
}

type RoleListQuery struct {
	Name   string `form:"name" filter:"like"`
	Status uint8  `form:"status" filter:"equal"`
}

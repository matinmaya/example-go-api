package usermodel

import (
	"reapp/internal/modules/user/rolemodel"
	"reapp/pkg/base/basemodel"
	"reapp/pkg/validators"
)

type User struct {
	basemodel.PrimaryKey
	Username string           `json:"username" gorm:"unique;not null;type:varchar(50);" validate:"required,min=6,max=50,unique=sys_users?id"`
	Password string           `json:"-" gorm:"not null;type:varchar(120);"`
	Status   bool             `json:"status" gorm:"not null;default=false;"`
	Roles    []rolemodel.Role `json:"roles,omitempty" gorm:"many2many:sys_user_role;"`
	RoleIds  []uint16         `json:"role_ids,omitempty" gorm:"-" validate:"required"`
	basemodel.SoftFields
	validators.ValidateScopeUnique
}

func (User) TableName() string {
	return "sys_users"
}

type UserRole struct {
	UserID    uint32                   `json:"user_id" gorm:"primaryKey;column:user_id;not null"`
	RoleID    uint16                   `json:"role_id" gorm:"primaryKey;column:role_id;not null"`
	User      User                     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role      rolemodel.Role           `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt basemodel.DateTimeFormat `json:"created_at"`
}

func (UserRole) TableName() string {
	return "sys_user_role"
}

type ChangePassword struct {
	UserID      uint32 `json:"user_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UserListQuery struct {
	Username string `form:"username" filter:"like"`
	Status   uint8  `form:"status" filter:"equal"`
}

package usermodel

import (
	"reapp/internal/modules/user/rolemodel"
	"reapp/pkg/base/basemodel"
	"reapp/pkg/validators"
)

type User struct {
	basemodel.PrimaryKey
	Username basemodel.TString `json:"username" gorm:"unique;not null;type:varchar(50);" validate:"required,min=6,max=50,slug_strict,unique=sys_users?id"`
	Password string            `json:"-" gorm:"not null;type:varchar(120);"`
	Status   bool              `json:"status" gorm:"not null;default=false;"`
	Img      string            `json:"img" gorm:"type:varchar(255);"`
	Roles    []rolemodel.Role  `json:"roles,omitempty" gorm:"many2many:sys_user_role;"`
	RoleIds  []uint16          `json:"role_ids,omitempty" gorm:"-" validate:"required"`
	basemodel.SoftFields
	validators.ValidateUniqueScope
}

func (User) TableName() string {
	return "sys_users"
}

type UserRole struct {
	UserID    uint32              `json:"user_id" gorm:"primaryKey;column:user_id;not null"`
	RoleID    uint16              `json:"role_id" gorm:"primaryKey;column:role_id;not null"`
	User      User                `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Role      rolemodel.Role      `gorm:"foreignKey:RoleID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	CreatedAt basemodel.TDateTime `json:"created_at"`
}

func (UserRole) TableName() string {
	return "sys_user_role"
}

type ChangePassword struct {
	UserID      uint32 `json:"user_id" validate:"required"`
	NewPassword string `json:"new_password" validate:"required,min=6"`
}

type UserListQuery struct {
	Username basemodel.TString `form:"username" filter:"like"`
	Status   uint8             `form:"status" filter:"equal"`
}

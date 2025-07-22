package usermodel

import "reapp/pkg/basemodel"

type TokenInfo struct {
	ID           uint64                   `json:"id" gorm:"primaryKey"`
	UserID       uint32                   `json:"user_id" gorm:"not null;index"`
	JTI          string                   `json:"jti" gorm:"not null;type:varchar(64);uniqueIndex"`
	RefreshToken string                   `json:"refresh_token" gorm:"not null;type:varchar(255);uniqueIndex"`
	Device       string                   `json:"device" gorm:"not null;type:varchar(100);"`
	Platform     string                   `json:"platform" gorm:"not null;type:varchar(50);"`
	Browser      string                   `json:"browser" gorm:"not null;type:varchar(50);"`
	OS           string                   `json:"os" gorm:"not null;type:varchar(50);"`
	Location     string                   `json:"location" gorm:"type:varchar(120);"`
	UserAgent    string                   `json:"user_agent" gorm:"not null;type:varchar(255);"`
	IP           string                   `json:"ip" gorm:"not null;type:varchar(45);"`
	ExpiresAt    basemodel.DateTimeFormat `json:"expires_at" gorm:"not null;index"`
	CreatedAt    basemodel.DateTimeFormat `json:"created_at"`
	UpdatedAt    basemodel.DateTimeFormat `json:"updated_at"`
}

func (TokenInfo) TableName() string {
	return "sys_token_infos"
}

type AuthChangePassword struct {
	OldPassword     string `json:"old_password" validate:"required,min=6"`
	NewPassword     string `json:"new_password" validate:"required,min=6"`
	ConfirmPassword string `json:"confirm_password" validate:"required,min=6"`
}

type AuthCredentials struct {
	Username string `json:"username" validate:"required,min=6"`
	Password string `json:"password" validate:"required,min=6"`
}

type AuthLoginResource struct {
	Username     string `json:"username"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

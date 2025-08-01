package basemodel

import (
	"gorm.io/gorm"
)

type MiniPrimaryKey struct {
	ID uint8 `json:"id" gorm:"primaryKey"`
}

type SmallPrimaryKey struct {
	ID uint16 `json:"id" gorm:"primaryKey"`
}

type PrimaryKey struct {
	ID uint32 `json:"id" gorm:"primaryKey"`
}

type BigPrimaryKey struct {
	ID uint64 `json:"id" gorm:"primaryKey"`
}

type SoftFields struct {
	CreatedBy *uint32        `json:"created_by"`
	UpdatedBy *uint32        `json:"updated_by"`
	DeletedBy *uint32        `json:"-"`
	CreatedAt DateTimeFormat `json:"created_at"`
	UpdatedAt DateTimeFormat `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

var bindingGuarded = []string{
	"id", "created_by", "updated_by", "deleted_by", "created_at", "updated_at", "deleted_at",
}

func (s *SoftFields) BindingGuarded() []string {
	return bindingGuarded
}

package basemodel

import (
	"encoding/json"

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
	CreatedAt TDateTime      `json:"created_at"`
	UpdatedAt TDateTime      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `json:"-" gorm:"index"`
}

type HttpLog struct {
	ID         uint64 `json:"id" gorm:"primaryKey"`
	Method     string `gorm:"size:10"`
	Path       string
	Query      string
	Body       string
	UserAgent  string
	IP         string
	StatusCode int
	CreatedAt  TDateTime `json:"created_at"`
}

func (HttpLog) TableName() string {
	return "sys_http_logs"
}

type TableLog struct {
	BigPrimaryKey
	TbName      string          `json:"tb_name" gorm:"type:varchar(50)"`
	TbID        uint64          `json:"tb_id"`
	Action      string          `json:"action" gorm:"type:varchar(10)"`
	ChangedData json.RawMessage `json:"changed_data"`
	FullData    json.RawMessage `json:"full_data"`
	CreatedBy   *uint32         `json:"created_by"`
	CreatedAt   TDateTime       `json:"created_at"`
}

func (TableLog) TableName() string {
	return "sys_tb_logs"
}

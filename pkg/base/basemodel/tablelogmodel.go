package basemodel

import "encoding/json"

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

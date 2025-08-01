package basemodel

type RequestLog struct {
	ID         uint64 `json:"id" gorm:"primaryKey"`
	Method     string `gorm:"size:10"`
	Path       string
	Query      string
	Body       string
	UserAgent  string
	IP         string
	StatusCode int
	CreatedAt  DateTimeFormat `json:"created_at"`
}

func (RequestLog) TableName() string {
	return "sys_request_logs"
}

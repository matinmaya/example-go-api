package customer

import (
	"reapp/pkg/base/basemodel"
	"reapp/pkg/validators"
)

type Customer struct {
	basemodel.PrimaryKey
	CardNumber basemodel.TString `json:"card_number" gorm:"unique;not null;type:varchar(16);" validate:"required,min=6,max=50,unique=cus_customers?id"`
	Fullname   basemodel.TString `json:"fullname" gorm:"not null;type:varchar(50);" validate:"required,min=6,max=50,unique=cus_customers?id"`
	Status     bool              `json:"status" gorm:"not null;default=false;"`
	Img        string            `json:"img" gorm:"type:varchar(255);"`
	basemodel.SoftFields
	validators.ValidateUniqueScope
}

func (Customer) TableName() string {
	return "cus_customers"
}

func (c Customer) GetID() uint32 {
	return c.ID
}

type CustomerListQuery struct {
	Fullname basemodel.TString `form:"fullname" filter:"like"`
	Status   bool              `form:"status" filter:"equal"`
}

package customer

import (
	"reapp/pkg/base/basemodel"
	"reapp/pkg/validators"
)

type Customer struct {
	basemodel.PrimaryKey
	Fullname       basemodel.TString   `json:"fullname" gorm:"not null;type:varchar(50);" validate:"required,min=6,max=50"`
	IdentifyNumber basemodel.TString   `json:"identify_number" gorm:"type:varchar(16);uniqueIndex;default:null" validate:"max=16,unique=cus_customers?nullable"`
	PassportNumber basemodel.TString   `json:"passport_number" gorm:"type:varchar(16);uniqueIndex;default:null" validate:"max=16,unique=cus_customers?nullable"`
	Phone          basemodel.TString   `json:"phone" gorm:"type:varchar(20);uniqueIndex;default:null" validate:"max=20,unique=cus_customers?nullable"`
	Email          basemodel.TString   `json:"email" gorm:"type:varchar(50);uniqueIndex;default:null" validate:"max=50,unique=cus_customers?nullable"`
	Gender         basemodel.TString   `json:"gender" gorm:"type:enum('M','F');default:null" validate:"omitempty,oneof=M F"`
	DOB            basemodel.TDateOnly `json:"dob" gorm:"type:date;default:null" validate:"omitempty,date"`
	DOBAddress     basemodel.TString   `json:"dob_address" gorm:"type:varchar(255);default:null"`
	Address        basemodel.TString   `json:"address" gorm:"type:varchar(255);"`
	Img            string              `json:"img" gorm:"type:varchar(255);"`
	Note           string              `json:"note" gorm:"type:text;"`
	Status         bool                `json:"status" gorm:"not null;default=false;"`
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
	Fullname       basemodel.TString `form:"fullname" filter:"like"`
	IdentifyNumber basemodel.TString `form:"identify_number" filter:"equal"`
	Status         bool              `form:"status" filter:"equal"`
}

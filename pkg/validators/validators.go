package validators

import (
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

var DB *gorm.DB
var Validate *validator.Validate

func InitValidation(db *gorm.DB, vlt *validator.Validate) *validator.Validate {
	DB = db
	Validate = vlt
	return Validate
}

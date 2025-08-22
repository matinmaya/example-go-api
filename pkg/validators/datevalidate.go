package validators

import (
	"time"

	"github.com/go-playground/validator/v10"
)

func Date(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	if fieldValue == "" {
		return true
	}

	_, err := time.Parse("2006-01-02", fieldValue)
	return err == nil
}

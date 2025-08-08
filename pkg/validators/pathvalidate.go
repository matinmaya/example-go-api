package validators

import (
	"fmt"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
)

func Path(fl validator.FieldLevel) bool {
	fieldValue := fl.Field().String()
	params := strings.Split(fl.Param(), "?")
	if len(params) < 1 {
		fmt.Printf("Messing directory in path validation")
		return false
	}

	_, err := os.Stat(fieldValue)

	return !os.IsNotExist(err)
}

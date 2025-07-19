package validators

import (
	"fmt"

	"github.com/go-playground/validator"
)

func GetMessage(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return "This field is required."
	case "min":
		return "Minimum length should be " + fe.Param()
	case "max":
		return "Maximum length should be " + fe.Param()
	case "gt":
		return "Value must be greater than " + fe.Param()
	case "lt":
		return "Value must be less than " + fe.Param()
	case "email":
		return "Invalid email format."
	case "numeric":
		return "Must be a numeric value."
	case "uuid":
		return "Invalid UUID format."
	case "unique":
		return "This value is already taken."
	case "path":
		return "Path does not exist."
	default:
		return fmt.Sprintf("Validation failed on field '%s' with tag '%s'", fe.Field(), fe.Tag())
	}
}

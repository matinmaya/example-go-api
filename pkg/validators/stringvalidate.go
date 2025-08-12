package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var slugStrictRegex = regexp.MustCompile(`^[A-Za-z0-9](?:[A-Za-z0-9_]{1,18}[A-Za-z0-9])$`)

func SlugStrict(fl validator.FieldLevel) bool {
	return slugStrictRegex.MatchString(fl.Field().String())
}

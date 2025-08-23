package reqvalidate

import (
	"encoding/json"
	"errors"
	"net/http"
	"reflect"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"reapp/pkg/http/response"
	"reapp/pkg/lang"
	"reapp/pkg/validators"
)

type IRequestDTO interface{}

func Validate(ctx *gin.Context, dto IRequestDTO) bool {
	bindingErr := ctx.ShouldBindJSON(dto)

	if err := validators.Validate.Struct(dto); err != nil {
		var validationErrors validator.ValidationErrors
		if errors.As(err, &validationErrors) {
			fieldErrors := make(map[string]string)

			reflectedType := reflect.TypeOf(dto).Elem()
			for _, fieldError := range validationErrors {
				fieldName := fieldError.Field()

				if field, found := reflectedType.FieldByName(fieldName); found {
					jsonTag := field.Tag.Get("json")
					if jsonTag != "" {
						jsonKey := strings.Split(jsonTag, ",")[0]
						fieldErrors[jsonKey] = validators.Message(ctx, fieldError)
					} else {
						fieldErrors[fieldName] = validators.Message(ctx, fieldError)
					}
				}
			}

			response.Error(ctx, http.StatusBadRequest, lang.Tran(ctx, "validation", "failed"), fieldErrors)
			return false
		}

		response.Error(ctx, http.StatusBadRequest, lang.Tran(ctx, "validation", "failed"), map[string]string{"error": err.Error()})
		return false
	}

	if bindingErr != nil {
		var unmarshalTypeErr *json.UnmarshalTypeError
		if errors.As(bindingErr, &unmarshalTypeErr) && unmarshalTypeErr.Type.Kind() == reflect.Bool {
			fieldName := unmarshalTypeErr.Field
			response.Error(ctx, http.StatusBadRequest, lang.Tran(ctx, "validation", "failed"), map[string]string{
				fieldName: lang.Tran(ctx, "validation", "boolean"),
			})
			return false
		}

		fieldName := unmarshalTypeErr.Field
		if fieldName != "" {
			response.Error(ctx, http.StatusBadRequest, lang.Tran(ctx, "validation", "failed"), map[string]string{
				fieldName: lang.Tran(ctx, "validation", "invalid_format"),
			})
			return false
		}

		response.Error(ctx, http.StatusBadRequest, lang.Tran(ctx, "validation", "failed"), map[string]string{"error": bindingErr.Error()})
		return false
	}

	return true
}

func ValidateParamID(ctx *gin.Context, id *uint64) bool {
	val, err := strconv.Atoi(ctx.Param("id"))
	if err != nil {
		response.Error(ctx, http.StatusBadRequest, err.Error(), nil)
		return false
	}

	*id = uint64(val)
	return true
}

package validators

import (
	"fmt"
	"reapp/pkg/lang"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator"
)

func GetMessage(ctx *gin.Context, fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return lang.Tran(ctx, "validation", "required")
	case "min":
		return lang.TranWithParams(ctx, "validation", "min", fe.Param())
	case "max":
		return lang.TranWithParams(ctx, "validation", "max", fe.Param())
	case "gt":
		return lang.TranWithParams(ctx, "validation", "gt", fe.Param())
	case "lt":
		return lang.TranWithParams(ctx, "validation", "lt", fe.Param())
	case "email":
		return lang.Tran(ctx, "validation", "email")
	case "numeric":
		return lang.Tran(ctx, "validation", "numeric")
	case "uuid":
		return lang.Tran(ctx, "validation", "uuid")
	case "unique":
		return lang.Tran(ctx, "validation", "unique")
	case "path":
		return lang.Tran(ctx, "validation", "path")
	default:
		return fmt.Sprintf("Validation failed on field '%s' with tag '%s'", fe.Field(), fe.Tag())
	}
}

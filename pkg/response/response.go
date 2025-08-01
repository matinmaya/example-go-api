package response

import (
	"net/http"
	"reapp/pkg/lang"

	"github.com/gin-gonic/gin"
)

type IResponse struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    interface{}       `json:"data,omitempty"`
	Errors  map[string]string `json:"errors,omitempty"`
}

func Success(ctx *gin.Context, code int, message string, data interface{}) {
	res := IResponse{
		Code:    code,
		Message: message,
	}
	if data != nil {
		res.Data = data
	}
	ctx.JSON(http.StatusOK, res)
}

func Error(ctx *gin.Context, code int, message string, err map[string]string) {
	res := IResponse{
		Code:    code,
		Message: message,
	}

	if err != nil {
		res.Errors = err
	}
	ctx.JSON(http.StatusOK, res)
}

func AsJSON(ctx *gin.Context, data interface{}, err interface{}) {
	if err != nil {
		if errMap, ok := err.(map[string]string); ok {
			Error(ctx, http.StatusInternalServerError, lang.ErrorMessage(ctx), errMap)
		} else {
			if errObj, ok := err.(error); ok {
				Error(ctx, http.StatusInternalServerError, lang.ErrorMessage(ctx), map[string]string{"error": errObj.Error()})
			} else {
				Error(ctx, http.StatusInternalServerError, lang.ErrorMessage(ctx), nil)
			}
		}
		return
	}
	Success(ctx, http.StatusOK, lang.SuccessMessage(ctx), data)
}

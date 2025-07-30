package lang

import (
	"reapp/internal/helpers/ctxhelper"
	langdata "reapp/internal/lang/data"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var data = map[string]map[string]map[string]string{
	"en": langdata.EN,
	"km": langdata.KM,
	"zh": langdata.ZH,
}

func Get(lang, group, key string) string {
	if langGroup, ok := data[lang]; ok {
		if messages, ok := langGroup[group]; ok {
			if val, ok := messages[key]; ok {
				return val
			}
		}
	}

	if messages, ok := data["en"][group]; ok {
		if val, ok := messages[key]; ok {
			return val
		}
	}

	return key
}

func Tran(ctx *gin.Context, group, key string) string {
	lang := ctx.MustGet("lang").(string)
	return Get(lang, group, key)
}

func TranByDB(db *gorm.DB, group, key string) string {
	lang := ctxhelper.GetLangByDBContext(db)
	return Get(lang, group, key)
}

func SuccessMessage(ctx *gin.Context) string {
	lang := ctx.MustGet("lang").(string)
	return Get(lang, "response", "success")
}

func ErrorMessage(ctx *gin.Context) string {
	lang := ctx.MustGet("lang").(string)
	return Get(lang, "response", "error")
}

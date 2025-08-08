package dbctx

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type TContextKey string

const dbContextKey = "db"
const langContextKey TContextKey = "lang"

func SetDBContext(ctx *gin.Context, db *gorm.DB) {
	ctxValue := context.WithValue(ctx.Request.Context(), langContextKey, ctx.MustGet("lang").(string))
	ctx.Request = ctx.Request.WithContext(ctxValue)
	ctx.Set(dbContextKey, db.WithContext(ctx.Request.Context()))
}

func DB(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet(dbContextKey).(*gorm.DB)
}

func DBContextKey() string {
	return dbContextKey
}

func LangByDBContext(db *gorm.DB) string {
	if db == nil {
		return "en"
	}
	if lang, ok := db.Statement.Context.Value(langContextKey).(string); ok {
		return lang
	}

	return "en"
}

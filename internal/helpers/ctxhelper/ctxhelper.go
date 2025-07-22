package ctxhelper

import (
	"context"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ctxKey string

const dbContextKey = "db"
const langContextKey ctxKey = "lang"

func SetDBContext(ctx *gin.Context, db *gorm.DB) {
	requestContext := ctx.Request.Context()
	requestContext = context.WithValue(requestContext, langContextKey, ctx.MustGet("lang").(string))
	dbContext := db.Session(&gorm.Session{Context: requestContext})
	ctx.Set(dbContextKey, dbContext)
}

func GetDB(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet(dbContextKey).(*gorm.DB)
}

func GetDBContextKey() string {
	return dbContextKey
}

func GetLangByDBContext(db *gorm.DB) string {
	if db == nil {
		return "en"
	}
	if lang, ok := db.Statement.Context.Value(langContextKey).(string); ok {
		return lang
	}
	return "en"
}

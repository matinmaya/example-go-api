package ctxhelper

import (
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

const ctxDBKey = "db"

func SetDBContext(ctx *gin.Context, db *gorm.DB) {
	reqCtx := ctx.Request.Context()
	ctxDB := db.Session(&gorm.Session{Context: reqCtx})
	ctx.Set(ctxDBKey, ctxDB)
}

func GetDB(ctx *gin.Context) *gorm.DB {
	return ctx.MustGet(ctxDBKey).(*gorm.DB)
}

func GetCxtDBKey() string {
	return ctxDBKey
}

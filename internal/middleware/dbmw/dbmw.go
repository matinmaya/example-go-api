package dbmw

import (
	"reapp/pkg/context/dbctx"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WithDBContext(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		dbctx.SetDBContext(ctx, db)
		ctx.Next()
	}
}

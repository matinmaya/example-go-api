package dbmdw

import (
	"reapp/pkg/helpers/ctxhelper"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func WithDBContext(db *gorm.DB) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctxhelper.SetDBContext(ctx, db)
		ctx.Next()
	}
}

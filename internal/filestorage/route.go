package filestorage

import (
	"github.com/gin-gonic/gin"
)

func RegisterFileRoute(r *gin.RouterGroup, storage Storage) {
	s := NewService(storage)
	fgh := NewFileHandler(s)
	r.POST("/uploads/file", fgh.Upload)
	r.POST("/uploads/image", fgh.UploadImage)
	r.GET("/*path", fgh.ServeFile)
	r.DELETE("/*path", fgh.Delete)
}

package filestorage

import (
	"github.com/gin-gonic/gin"
)

func RegisterFileRoute(r *gin.RouterGroup, storage Storage) {
	s := NewService(storage)
	fgh := NewFileHandler(s)
	r.POST("/upload", fgh.Upload)
	r.POST("/image/upload", fgh.UploadImage)
	r.GET("/*path", fgh.ServeFile)
	r.DELETE("/*path", fgh.Delete)
}

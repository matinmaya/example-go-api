package loggermdw

import (
	"bytes"
	"io"
	"reapp/pkg/base/basemodel"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func RequestLogger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		log := basemodel.RequestLog{
			Method:     c.Request.Method,
			Path:       c.Request.URL.Path,
			Query:      c.Request.URL.RawQuery,
			Body:       truncate(string(bodyBytes), 1000),
			UserAgent:  c.Request.UserAgent(),
			IP:         c.ClientIP(),
			StatusCode: c.Writer.Status(),
		}

		go db.Create(&log)
	}
}

func truncate(s string, max int) string {
	if len(s) > max {
		return s[:max] + "..."
	}
	return s
}

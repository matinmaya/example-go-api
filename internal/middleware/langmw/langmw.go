package langmw

import (
	"strings"

	"github.com/gin-gonic/gin"
)

func Language() gin.HandlerFunc {
	return func(c *gin.Context) {
		lang := c.Query("lang")
		if lang == "" {
			acceptLang := c.GetHeader("Accept-Language")
			if len(acceptLang) >= 2 {
				lang = strings.ToLower(acceptLang[:2])
			}
		}

		if lang != "en" && lang != "km" && lang != "zh" {
			lang = "en"
		}

		c.Set("lang", lang)
		c.Next()
	}
}

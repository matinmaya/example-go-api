package bucket

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"reapp/pkg/redisclient"
)

func GetImagePath(fullURL string) string {
	u, err := url.Parse(fullURL)
	if err != nil {
		return fullURL
	}
	return u.Path
}

func IsFullImagePath(path string) bool {
	lower := strings.ToLower(path)
	if strings.HasPrefix(lower, "/") {
		return false
	}
	return strings.HasPrefix(lower, "http://") || strings.HasPrefix(lower, "https://")
}

func IsAbsoluteImagePath(path string) bool {
	lower := strings.ToLower(path)
	return strings.HasPrefix(lower, "/")
}

func GetFullImageURL(ctx *gin.Context, imagePath string) string {
	redisClient := redisclient.Client()
	baseURL, _ := redisClient.Get("bucket_base_url").Result()
	if baseURL == "" {
		baseURL = getCurrentServerURL(ctx)
	}
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(imagePath, "/")
}

func getCurrentServerURL(ctx *gin.Context) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + ctx.Request.Host + "/storages"
}

package filesystem

import (
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"

	"reapp/pkg/redisclient"
)

var prefixRoutePath string

func currentServerURL(ctx *gin.Context) string {
	scheme := "http"
	if ctx.Request.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + ctx.Request.Host + "/" + TrimPath(prefixRoutePath)
}

func SetPrefixRoutePath(prefixPath string) {
	prefixRoutePath = prefixPath
}

func PrefixRoutePath() string {
	return prefixRoutePath
}

func ImagePath(fullURL string) string {
	u, err := url.Parse(fullURL)
	if err != nil {
		return fullURL
	}
	return u.Path
}

func FullImageURL(ctx *gin.Context, imagePath string) string {
	redisClient := redisclient.Client()
	baseURL, _ := redisClient.Get("bucket_base_url").Result()
	if baseURL == "" {
		baseURL = currentServerURL(ctx)
	}
	return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(imagePath, "/")
}

func TrimPath(path string) string {
	str := strings.TrimRight(path, "/")
	return strings.TrimLeft(str, "/")
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

func IsValidStoragePath(path string) bool {
	if strings.HasPrefix(path, "/") {
		return false
	}
	trimmed := strings.TrimPrefix(path, "./")
	segments := strings.Split(trimmed, "/")
	return len(segments) == 2 && segments[0] != "" && segments[1] != ""
}

package redishelper

import (
	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var repoCacheDuration time.Duration

func InitRedis(client *redis.Client, cacheTTL int) {
	redisClient = client
	repoCacheDuration = time.Duration(cacheTTL) * time.Minute
}

func Client() *redis.Client {
	if redisClient == nil {
		panic("Redis client is not initialized")
	}
	return redisClient
}

func GetRepoCacheDuration() time.Duration {
	return repoCacheDuration
}

package redisclient

import (
	"time"

	"github.com/go-redis/redis"
)

var redisClient *redis.Client
var repoCacheDur time.Duration

func InitRedis(client *redis.Client, repoCacheTTL int) {
	redisClient = client
	repoCacheDur = time.Duration(repoCacheTTL) * time.Minute
}

func Client() *redis.Client {
	if redisClient == nil {
		panic("Redis client is not initialized")
	}
	return redisClient
}

func RepoCacheDur() time.Duration {
	return repoCacheDur
}

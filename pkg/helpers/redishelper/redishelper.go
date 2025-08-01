package redishelper

import "github.com/go-redis/redis"

var redisClient *redis.Client

func InitRedis(client *redis.Client) {
	redisClient = client
}

func Client() *redis.Client {
	if redisClient == nil {
		panic("Redis client is not initialized")
	}
	return redisClient
}

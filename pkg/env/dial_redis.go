package env

import (
	"fmt"

	"github.com/go-redis/redis"
)

func DialRedis(cf *Config) (*redis.Client, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", cf.Redis.Host, cf.Redis.Port),
		Password: cf.Redis.Password,
		DB:       cf.Redis.DB,
	})

	if err := client.Ping().Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	fmt.Println("Connected to Redis")
	return client, nil
}

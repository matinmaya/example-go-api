package rediservice

import (
	"encoding/json"
	"fmt"
	"log"

	"reapp/pkg/redisclient"
)

func repositoryRoot() string {
	return "repositories"
}

func CacheOfRepository[T any](namespace string, collection string, key string, data *T) error {
	validate := map[string]string{
		"namespace":  namespace,
		"collection": collection,
		"key":        key,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	cacheKey := fmt.Sprintf("%s:%s:%s:%s", repositoryRoot(), namespace, collection, key)
	redisClient := redisclient.Client()
	cached, err := redisClient.Get(cacheKey).Result()
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	if err := json.Unmarshal([]byte(cached), data); err != nil {
		return err
	}
	return nil
}

func SetCacheOfRepository[T any](namespace string, collection string, key string, params T) error {
	validate := map[string]string{
		"namespace":  namespace,
		"collection": collection,
		"key":        key,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	cacheKey := fmt.Sprintf("%s:%s:%s:%s", repositoryRoot(), namespace, collection, key)
	redisClient := redisclient.Client()
	data, err := json.Marshal(params)
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	if err := redisClient.Set(cacheKey, data, redisclient.RepoCacheDur()).Err(); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	return nil
}

func ClearCacheOfRepository(namespace string) error {
	validate := map[string]string{
		"namespace": namespace,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	prefixKey := fmt.Sprintf("%s:%s:*", repositoryRoot(), namespace)
	redisClient := redisclient.Client()
	iter := redisClient.Scan(0, prefixKey, 0).Iterator()
	for iter.Next() {
		if err := redisClient.Del(iter.Val()).Err(); err != nil {
			log.Printf("%s", err.Error())
			return err
		}
	}

	return nil
}

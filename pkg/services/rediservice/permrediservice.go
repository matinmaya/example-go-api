package rediservice

import (
	"encoding/json"
	"fmt"
	"log"
	"reapp/pkg/redisclient"
	"time"
)

func CacheOfPerms(userID string) ([]string, error) {
	validate := map[string]string{
		"userID": userID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return nil, err
	}
	cacheKey := "permissions:user:" + userID
	redisClient := redisclient.Client()
	cached, err := redisClient.Get(cacheKey).Result()
	if err != nil {
		return nil, err
	}
	if cached == "" {
		return nil, fmt.Errorf("no permissions found for userID: %s", userID)
	}

	var permissions []string
	if err := json.Unmarshal([]byte(cached), &permissions); err != nil {
		log.Printf("%s", err.Error())
		return nil, err
	}
	return permissions, nil
}

func SetCacheOfPerms(userID string, permissions []string) error {
	validate := map[string]string{
		"userID": userID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	cacheKey := "permissions:user:" + userID
	redisClient := redisclient.Client()
	data, err := json.Marshal(permissions)
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	if err := redisClient.Set(cacheKey, data, time.Hour).Err(); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	return nil
}

func ClearCacheOfPerms() error {
	redisClient := redisclient.Client()
	iter := redisClient.Scan(0, "permissions:user:*", 0).Iterator()
	for iter.Next() {
		if err := redisClient.Del(iter.Val()).Err(); err != nil {
			log.Printf("%s", err.Error())
			return err
		}
	}
	return nil
}

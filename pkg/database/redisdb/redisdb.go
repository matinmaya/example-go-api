package redisdb

import (
	"encoding/json"
	"fmt"
	"reapp/pkg/helpers/redishelper"
	"time"
)

func GetCacheOfPerms(userID string) ([]string, error) {
	if userID == "" {
		return nil, fmt.Errorf("userID cannot be empty")
	}
	cacheKey := "permissions:user:" + userID
	redisClient := redishelper.Client()
	cached, err := redisClient.Get(cacheKey).Result()
	if err != nil {
		return nil, err
	}
	if cached == "" {
		return nil, fmt.Errorf("no permissions found for userID: %s", userID)
	}

	var permissions []string
	if err := json.Unmarshal([]byte(cached), &permissions); err != nil {
		return nil, err
	}
	return permissions, nil
}

func SetCacheOfPerms(userID string, permissions []string) error {
	cacheKey := "permissions:user:" + userID
	redisClient := redishelper.Client()
	data, err := json.Marshal(permissions)
	if err != nil {
		return err
	}

	if err := redisClient.Set(cacheKey, data, time.Hour).Err(); err != nil {
		return err
	}
	return nil
}

func ClearCacheOfPerms() error {
	redisClient := redishelper.Client()
	iter := redisClient.Scan(0, "permissions:user:*", 0).Iterator()
	for iter.Next() {
		if err := redisClient.Del(iter.Val()).Err(); err != nil {
			return err
		}
	}
	return nil
}

func GetCacheOfAuthUser[T any](userID string, auth *T) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	cacheKey := "auth:user:" + userID
	redisClient := redishelper.Client()
	cached, err := redisClient.Get(cacheKey).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(cached), auth); err != nil {
		return err
	}
	return nil
}

func SetCacheOfAuthUser[T any](userID string, auth T) error {
	if userID == "" {
		return fmt.Errorf("userID cannot be empty")
	}
	cacheKey := "auth:user:" + userID
	redisClient := redishelper.Client()
	data, err := json.Marshal(auth)
	if err != nil {
		return err
	}
	if err := redisClient.Set(cacheKey, data, 10*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}

func RemoveCacheOfAuthUser(userID string) error {
	if userID == "" {
		return nil
	}
	cacheKey := "auth:user:" + userID
	redisClient := redishelper.Client()
	if err := redisClient.Del(cacheKey).Err(); err != nil {
		return err
	}
	return nil
}

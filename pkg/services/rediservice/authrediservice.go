package rediservice

import (
	"encoding/json"
	"log"
	"reapp/pkg/redisclient"
	"time"
)

func CacheOfAuthUser[T any](userID string, auth *T) error {
	validate := map[string]string{
		"userID": userID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	cacheKey := "auth:user:" + userID
	redisClient := redisclient.Client()
	cached, err := redisClient.Get(cacheKey).Result()
	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(cached), auth); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	return nil
}

func SetCacheOfAuthUser[T any](userID string, auth T) error {
	validate := map[string]string{
		"userID": userID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	cacheKey := "auth:user:" + userID
	redisClient := redisclient.Client()
	data, err := json.Marshal(auth)
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	if err := redisClient.Set(cacheKey, data, 10*time.Minute).Err(); err != nil {
		return err
	}
	return nil
}

func RemoveCacheOfAuthUser(userID string) error {
	validate := map[string]string{
		"userID": userID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	cacheKey := "auth:user:" + userID
	redisClient := redisclient.Client()
	if err := redisClient.Del(cacheKey).Err(); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	return nil
}

func RevokeToken(tokenID string, expiresAt time.Time) error {
	validate := map[string]string{
		"tokenID": tokenID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	redisClient := redisclient.Client()
	cacheKey := "revoked:" + tokenID
	return redisClient.Set(cacheKey, "true", time.Until(expiresAt)).Err()
}

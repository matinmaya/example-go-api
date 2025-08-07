package redisdb

import (
	"encoding/json"
	"fmt"
	"log"
	"reapp/pkg/helpers/redishelper"
	"time"
)

// private func
func validateNotEmpty(params map[string]string) error {
	for name, value := range params {
		if value == "" {
			return fmt.Errorf("argument [%s] cannot be empty", name)
		}
	}
	return nil
}

func GetCacheOfPerms(userID string) ([]string, error) {
	validate := map[string]string{
		"userID": userID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return nil, err
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
	redisClient := redishelper.Client()
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
	redisClient := redishelper.Client()
	iter := redisClient.Scan(0, "permissions:user:*", 0).Iterator()
	for iter.Next() {
		if err := redisClient.Del(iter.Val()).Err(); err != nil {
			log.Printf("%s", err.Error())
			return err
		}
	}
	return nil
}

func GetCacheOfAuthUser[T any](userID string, auth *T) error {
	validate := map[string]string{
		"userID": userID,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	cacheKey := "auth:user:" + userID
	redisClient := redishelper.Client()
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
	redisClient := redishelper.Client()
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
	redisClient := redishelper.Client()
	if err := redisClient.Del(cacheKey).Err(); err != nil {
		log.Printf("%s", err.Error())
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

	cacheKey := fmt.Sprintf("repositories:%s:%s:%s", namespace, collection, key)
	redisClient := redishelper.Client()
	data, err := json.Marshal(params)
	if err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	if err := redisClient.Set(cacheKey, data, redishelper.GetRepoCacheDuration()).Err(); err != nil {
		log.Printf("%s", err.Error())
		return err
	}
	return nil
}

func GetCacheOfRepository[T any](namespace string, collection string, key string, data *T) error {
	validate := map[string]string{
		"namespace":  namespace,
		"collection": collection,
		"key":        key,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	cacheKey := fmt.Sprintf("repositories:%s:%s:%s", namespace, collection, key)
	redisClient := redishelper.Client()
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

func ClearCacheOfRepository(namespace string) error {
	validate := map[string]string{
		"namespace": namespace,
	}
	if err := validateNotEmpty(validate); err != nil {
		log.Printf("%s", err.Error())
		return err
	}

	prefixKey := fmt.Sprintf("repositories:%s:*", namespace)
	redisClient := redishelper.Client()
	iter := redisClient.Scan(0, prefixKey, 0).Iterator()
	for iter.Next() {
		if err := redisClient.Del(iter.Val()).Err(); err != nil {
			log.Printf("%s", err.Error())
			return err
		}
	}

	return nil
}

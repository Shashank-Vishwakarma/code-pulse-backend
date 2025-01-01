package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func InitializeRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", config.Config.REDIS_ENDPOINT, config.Config.REDIS_PORT),
		Password: config.Config.REDIS_PASSWORD,
		DB:       0, // use default DB
	})
}

func SetCache(key string, value string) error {
	ctx := context.Background()

	result := RedisClient.Set(ctx, key, value, time.Hour*24)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

func GetCache(key string) (string, error) {
	ctx := context.Background()

	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}

	return value, nil
}

func InvalidateCache(key string) error {
	ctx := context.Background()

	result := RedisClient.Del(ctx, key)
	if result.Err() != nil {
		return result.Err()
	}

	return nil
}

package services

import (
	"context"
	"fmt"
	"time"

	"github.com/Shashank-Vishwakarma/code-pulse-backend/pkg/config"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
)

var RedisClient *redis.Client

func InitializeRedis() {
	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%s", config.Config.REDIS_ENDPOINT, config.Config.REDIS_PORT),
		// Password:       config.Config.REDIS_PASSWORD,
		DB: 2,
	})

	err := client.Ping(context.Background()).Err()
	if err != nil {
		logrus.Infof("Failed to connect to Redis: %v", err)
		panic(err)
	}

	logrus.Println("Connected to Redis successfully")

	RedisClient = client
}

func SetCache(key string, value interface{}, expiration time.Duration) error {
	ctx := context.Background()

	err := RedisClient.Set(ctx, key, value, expiration).Err()
	if err != nil {
		return err
	}

	return nil
}

func GetCache(key string) interface{} {
	ctx := context.Background()

	value, err := RedisClient.Get(ctx, key).Result()
	if err != nil {
		return nil
	}

	return value
}

func InvalidateCache(key string) error {
	ctx := context.Background()

	err := RedisClient.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	return nil
}

package services

import (
	"fmt"

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

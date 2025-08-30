package config

import (
	"context"
	"log"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var RedisClient *redis.Client

func InitRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr: viper.GetString("redis.addr"),
		DB:   viper.GetInt("redis.db"),
	})

	ctx := context.Background()
	_, err := RedisClient.Ping(ctx).Result()
	if err != nil {
		log.Fatal("Failed to connect to Redis:", err)
	}
	log.Println("Redis connected successfully")
}

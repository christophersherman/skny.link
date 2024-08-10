package config

import (
	"context"
	"log"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

func NewRedisClient() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "redis:6379", // Use the correct Redis server address in production
		Password: "",           // No password set if empty
		DB:       0,            // Default DB
	})

	// Ping the Redis server to ensure it's working
	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Fatalf("Could not connect to Redis: %v", err)
	}
	log.Printf("Connected to Redis: %s", pong)

	return rdb
}

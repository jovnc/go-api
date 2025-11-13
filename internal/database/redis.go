package database

import (
	"context"
	"go_api/internal/config"
	"log"
	"strconv"

	"github.com/redis/go-redis/v9"
)

var RedisClient *redis.Client

func ConnectRedis() *redis.Client {
	// Get Redis config from config
	addr := config.GlobalConfig.RedisAddr
	password := config.GlobalConfig.RedisPassword
	db := config.GlobalConfig.RedisDB
	dbInt, err := strconv.Atoi(db)
	if err != nil {
		log.Fatalf("failed to convert Redis DB to int: %v", err)
	}

	// Create Redis client
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       int(dbInt),
	})

	// Test connection
	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
		log.Fatalf("failed to ping Redis: %v", err)
	}

	log.Println("Redis connection established successfully")

	return RedisClient
}

func GetRedisClient() *redis.Client {
	return RedisClient
}

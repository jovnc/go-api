package storage

import (
	"context"
	"log"
	"strconv"

	"go_api/internal/config"

	"github.com/redis/go-redis/v9"
	"github.com/redis/go-redis/v9/maintnotifications"
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
		// Explicitly disable maintenance notifications
		// This prevents the client from sending CLIENT MAINT_NOTIFICATIONS ON
		MaintNotificationsConfig: &maintnotifications.Config{
			Mode: maintnotifications.ModeDisabled,
		},
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

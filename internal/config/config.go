package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort   string
	DatabaseURL  string
	Environment  string
	LogLevel     string
	JWTSecretKey string
	RedisAddr    string
	RedisPassword string
	RedisDB       string
}

var GlobalConfig *Config

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		return nil, fmt.Errorf("error loading file: %v", err)
	}

	databaseURL := getEnv("DATABASE_URL", "")
	if databaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}

	jwtSecretKey := getEnv("JWT_SECRET_KEY", "")
	if jwtSecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required")
	}

	GlobalConfig = &Config{
		ServerPort:   getEnv("SERVER_PORT", "8080"),
		DatabaseURL:  databaseURL,
		Environment:  getEnv("ENVIRONMENT", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
		JWTSecretKey: jwtSecretKey,
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnv("REDIS_DB", "0"),
	}

	return GlobalConfig, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

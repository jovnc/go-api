package config

import (
	"fmt"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	ServerPort        string
	DatabaseURL       string
	DatabaseURLPooler string
	Environment       string
	LogLevel          string
	JWTSecretKey      string
	RedisAddr         string
	RedisPassword     string
	RedisDB           string
	RateLimit         int
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

	databaseURLPooler := getEnv("DATABASE_URL_POOLER", "")
	if databaseURLPooler == "" {
		return nil, fmt.Errorf("DATABASE_URL_POOLER is required")
	}

	jwtSecretKey := getEnv("JWT_SECRET_KEY", "")
	if jwtSecretKey == "" {
		return nil, fmt.Errorf("JWT_SECRET_KEY is required")
	}

	rateLimit, err := strconv.Atoi(getEnv("RATE_LIMIT", "100"))
	if err != nil {
		return nil, fmt.Errorf("RATE_LIMIT is not a valid integer: %v", err)
	}

	GlobalConfig = &Config{
		ServerPort:        getEnv("SERVER_PORT", "8080"),
		DatabaseURL:       databaseURL,
		DatabaseURLPooler: databaseURLPooler,
		Environment:       getEnv("ENVIRONMENT", "development"),
		LogLevel:          getEnv("LOG_LEVEL", "info"),
		JWTSecretKey:      jwtSecretKey,
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           getEnv("REDIS_DB", "0"),
		RateLimit:         rateLimit,
	}

	return GlobalConfig, nil
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

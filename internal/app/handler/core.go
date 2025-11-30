package handler

import (
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	DB    *gorm.DB
	Redis *redis.Client
}

func NewHandler(db *gorm.DB, redis *redis.Client) *Handler {
	return &Handler{
		DB:    db,
		Redis: redis,
	}
}

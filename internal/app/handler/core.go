package handler

import (
	"go_api/internal/app/repository"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type Handler struct {
	DB    *gorm.DB
	BlogRepository *repository.BlogRepository
	Redis *redis.Client
}

func NewHandler(db *gorm.DB, blogRepository *repository.BlogRepository, redis *redis.Client) *Handler {
	return &Handler{
		DB:    db,
		BlogRepository: blogRepository,
		Redis: redis,
	}
}

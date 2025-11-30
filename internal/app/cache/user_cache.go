package cache

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"go_api/internal/app/model"
	"time"

	"github.com/redis/go-redis/v9"
)

type UserCache struct {
	redis *redis.Client
}

func NewUserCache(redis *redis.Client) *UserCache {
	return &UserCache{redis: redis}
}

func (c *UserCache) GetUser(ctx context.Context, id uint) (*model.User, error) {
	cacheKey := fmt.Sprintf("user:%d", id)
	if cached, err := c.redis.Get(ctx, cacheKey).Result(); err == nil {
		var user model.User
		if err := json.Unmarshal([]byte(cached), &user); err == nil {
			return &user, nil
		}
	}
	return nil, errors.New("user not found")
}

func (c *UserCache) SetUser(ctx context.Context, id uint, user *model.User) error {
	cacheKey := fmt.Sprintf("user:%d", id)
	userJSON, err := json.Marshal(user)
	if err != nil {
		return err
	}
	return c.redis.Set(ctx, cacheKey, userJSON, time.Minute*5).Err()
}

func (c *UserCache) BlacklistToken(ctx context.Context, token string) error {
	return c.redis.Set(ctx, token, "blacklisted", time.Minute*5).Err()
}

func (c *UserCache) CleanUserSession(ctx context.Context, userId uint) error {
	userIdStr := fmt.Sprintf("user:%d", userId)
	iter := c.redis.Scan(ctx, 0, userIdStr+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := c.redis.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("failed to clean user session: %v", err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan redis keys: %v", err)
	}
	return nil
}

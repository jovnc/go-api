package service

import (
	"context"
	"errors"
	"time"

	"go_api/internal/app/cache"
	"go_api/internal/app/dto"
	"go_api/internal/app/model"
	"go_api/internal/app/repository"
	"go_api/internal/config"
	"go_api/internal/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

var (
	ErrUserNotFound     = errors.New("user not found")
	ErrInvalidPassword  = errors.New("invalid password")
	ErrTokenGeneration  = errors.New("failed to generate token")
	ErrPasswordHashing  = errors.New("failed to hash password")
	ErrUserCreation     = errors.New("failed to create user")
	ErrCacheOperation   = errors.New("cache operation failed")
	ErrTokenBlacklist   = errors.New("failed to blacklist token")
	ErrSessionCleanup   = errors.New("failed to clean user session")
)

type UserService struct {
	repo  *repository.UserRepository
	cache *cache.UserCache
}

func NewUserService(db *gorm.DB, redis *redis.Client) *UserService {
	return &UserService{
		repo:  repository.NewUserRepository(db),
		cache: cache.NewUserCache(redis),
	}
}

// GetUserProfile retrieves user profile, checking cache first then database
func (s *UserService) GetUserProfile(ctx context.Context, userID uint) (*model.User, bool, error) {
	// Try cache first
	user, err := s.cache.GetUser(ctx, userID)
	if err == nil && user != nil {
		return user, true, nil
	}

	// Get from database
	user, err = s.repo.FindByID(ctx, userID)
	if err != nil {
		return nil, false, ErrUserNotFound
	}

	// Cache the user
	if err := s.cache.SetUser(ctx, userID, user); err != nil {
		return nil, false, ErrCacheOperation
	}

	return user, false, nil
}

// CreateUser creates a new user with hashed password
func (s *UserService) CreateUser(ctx context.Context, req dto.CreateUserRequest) (*model.User, error) {
	hashedPassword, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, ErrPasswordHashing
	}

	user := &model.User{
		Username: req.Username,
		Email:    req.Email,
		Password: hashedPassword,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, ErrUserCreation
	}

	return user, nil
}

// LoginUser authenticates user and returns a JWT token
func (s *UserService) LoginUser(ctx context.Context, req dto.LoginUserRequest) (string, error) {
	user, err := s.repo.FindByEmail(ctx, req.Email)
	if err != nil {
		return "", ErrUserNotFound
	}

	if !util.ComparePassword(req.Password, user.Password) {
		return "", ErrInvalidPassword
	}

	token, err := util.GenerateToken(user.ID, user.Username, []byte(config.GlobalConfig.JWTSecretKey))
	if err != nil {
		return "", ErrTokenGeneration
	}

	return token, nil
}

// LogoutUser blacklists the token and cleans user session
func (s *UserService) LogoutUser(ctx context.Context, userID uint, token string, expiresAt *jwt.NumericDate) error {
	// Calculate TTL - default to 5 minutes if expiration is not set
	ttl := time.Minute * 5
	if expiresAt != nil {
		ttl = time.Until(expiresAt.Time)
		// If token is expired, use default TTL
		if ttl <= 0 {
			ttl = time.Minute * 5
		}
	}

	// Blacklist token
	if err := s.cache.BlacklistToken(ctx, token); err != nil {
		return ErrTokenBlacklist
	}

	// Clean user session from cache
	if err := s.cache.CleanUserSession(ctx, userID); err != nil {
		return ErrSessionCleanup
	}

	return nil
}

// ListAllUsers retrieves all users from the database
func (s *UserService) ListAllUsers(ctx context.Context) ([]model.User, error) {
	users, err := s.repo.ListAllUsers(ctx)
	if err != nil {
		return nil, err
	}
	return users, nil
}

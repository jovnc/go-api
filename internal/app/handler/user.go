package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"go_api/internal/app/dto"
	"go_api/internal/app/model"
	"go_api/internal/config"
	"go_api/internal/middleware"
	"go_api/internal/util"
)

// UserProfileHandler returns the user profile
func (h *Handler) UserProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get claims from context
		claims, ok := r.Context().Value(middleware.UserClaimsKey).(*util.Claims)
		if !ok {
			util.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		// Get user from Redis (if exists)
		cacheKey := fmt.Sprintf("user:%d", claims.UserID)
		if cached, err := h.Redis.Get(ctx, cacheKey).Result(); err == nil {
			var user model.User
			if err := json.Unmarshal([]byte(cached), &user); err == nil {
				util.ResponseWithSuccess(w, http.StatusOK, "User profile (from cache)", user)
				return
			}
		}

		// Get user from database
		user := &model.User{}
		if err := h.DB.WithContext(ctx).Where("id = ?", claims.UserID).First(user).Error; err != nil {
			util.ResponseWithError(w, http.StatusNotFound, "User not found", err.Error())
			return
		}

		// Cache user in Redis
		userJSON, err := json.Marshal(user)
		if err == nil {
			h.Redis.Set(ctx, cacheKey, userJSON, time.Minute*5)
		}

		util.ResponseWithSuccess(w, http.StatusOK, "User profile (from database)", user)
	}
}

// CreateUserHandler creates a new user
func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Validate request
		if err := util.Validate(req); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Create user
		hashedPassword, err := util.HashPassword(req.Password)
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to hash password", err.Error())
			return
		}

		user := &model.User{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
		}

		if err := h.DB.WithContext(ctx).Create(user).Error; err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to create user", err.Error())
			return
		}

		util.ResponseWithSuccess(w, http.StatusCreated, "User created successfully", user)
	}
}

// LoginUserHandler logs in a user
func (h *Handler) LoginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.LoginUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Validate request
		if err := util.Validate(req); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Fetch user from DB
		user := &model.User{}
		if err := h.DB.WithContext(ctx).Where("email = ?", req.Email).First(user).Error; err != nil {
			util.ResponseWithError(w, http.StatusNotFound, "User not found", err.Error())
			return
		}

		// Compare password
		if !util.ComparePassword(req.Password, user.Password) {
			util.ResponseWithError(w, http.StatusUnauthorized, "Invalid password", "")
			return
		}

		// Generate token
		token, err := util.GenerateToken(user.ID, user.Username, []byte(config.GlobalConfig.JWTSecretKey))
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
			return
		}

		// Return token
		util.ResponseWithSuccess(w, http.StatusOK, "Login successful", map[string]string{
			"token": token,
		})
	}
}

// LogoutUserHandler logs out a user
func (h *Handler) LogoutUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get claims from context
		claims, ok := r.Context().Value(middleware.UserClaimsKey).(*util.Claims)
		if !ok {
			util.ResponseWithError(w, http.StatusBadRequest, "Missing claims", "Missing claims")
			return
		}

		// Extract token from Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.ResponseWithError(w, http.StatusUnauthorized, "Missing Authorization header", "Missing Authorization header")
			return
		}

		token, err := extractTokenFromHeader(authHeader)
		if err != nil {
			util.ResponseWithError(w, http.StatusUnauthorized, "Missing authorization token", err.Error())
			return
		}

		// Convert expireAt to time.Time
		expirationTime := time.Unix(claims.ExpiresAt, 0)
		now := time.Now()
		ttl := expirationTime.Sub(now)

		// If token is expired, set TTL to 5 minutes
		if ttl <= 0 {
			ttl = time.Minute * 5
		}

		// Blacklist token from Redis
		err = h.Redis.Set(ctx, token, "blacklisted", ttl).Err()
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to blacklist token", err.Error())
			return
		}

		// Clean user profile data from Redis
		userIdStr := fmt.Sprintf("user:%d", claims.UserID)
		if err := h.cleanUserSession(ctx, userIdStr); err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to clean user session", err.Error())
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "Logout successful", nil)
	}
}

// ListAllUsersHandler lists all users (with optional pagination)
func (h *Handler) ListAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		users := []model.User{}
		if err := h.DB.WithContext(ctx).Find(&users).Error; err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to list users", err.Error())
			return
		}
		util.ResponseWithSuccess(w, http.StatusOK, "List of all users", users)
	}
}

// Helper functions

// cleanUserSession cleans the user session from Redis
func (h *Handler) cleanUserSession(ctx context.Context, userIdStr string) error {
	iter := h.Redis.Scan(ctx, 0, userIdStr+"*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		if err := h.Redis.Del(ctx, key).Err(); err != nil {
			return fmt.Errorf("failed to clean user session: %v", err)
		}
	}
	if err := iter.Err(); err != nil {
		return fmt.Errorf("failed to scan redis keys: %v", err)
	}
	return nil
}

// extractTokenFromHeader extracts the token from the Authorization header
func extractTokenFromHeader(authHeader string) (string, error) {
	tokenString := strings.TrimPrefix(authHeader, "Bearer ")
	if tokenString == "" {
		return "", fmt.Errorf("missing authorization token")
	}
	return tokenString, nil
}

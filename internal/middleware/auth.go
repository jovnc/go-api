package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"go_api/internal/config"
	"go_api/internal/storage"
	"go_api/internal/util"
)

type contextKey string

const UserClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			util.ResponseWithError(w, http.StatusUnauthorized, "Missing Authorization header", "Authorization header is required")
			return
		}

		// Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &util.Claims{}

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GlobalConfig.JWTSecretKey), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				util.ResponseWithError(w, http.StatusUnauthorized, "Signature invalid", "Signature invalid")
			}
			util.ResponseWithError(w, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}

		// Validate token
		if !token.Valid {
			util.ResponseWithError(w, http.StatusUnauthorized, "Invalid token", "Token is invalid")
			return
		}

		// Check if token is blacklisted from Redis
		redisClient := storage.GetRedisClient()
		if redisClient == nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to get Redis client", "Failed to get Redis client")
			return
		}

		isBlacklisted, err := redisClient.Get(r.Context(), tokenString).Result()
		if err == nil && isBlacklisted == "blacklisted" {
			util.ResponseWithError(w, http.StatusUnauthorized, "Token is blacklisted", "Token is blacklisted")
			return
		}

		// Set context
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

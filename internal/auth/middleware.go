package auth

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt"

	"go_api/internal/config"
	"go_api/internal/utils"
)

type contextKey string

const UserClaimsKey contextKey = "claims"

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Get Authorization header
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Missing Authorization header", "Authorization header is required")
			return
		}

		// Bearer token
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		claims := &utils.Claims{}

		// Parse token
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return []byte(config.GlobalConfig.JWTSecretKey), nil
		})
		if err != nil {
			if errors.Is(err, jwt.ErrSignatureInvalid) {
				utils.ResponseWithError(w, http.StatusUnauthorized, "Signature invalid", "Signature invalid")
			}
			utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid token", err.Error())
			return
		}

		// Validate token
		if !token.Valid {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid token", "Token is invalid")
			return
		}

		// Set context
		ctx := context.WithValue(r.Context(), UserClaimsKey, claims)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}

package handler

import (
	"encoding/json"
	"errors"
	"net/http"

	"go_api/internal/app/dto"
	"go_api/internal/app/service"
	"go_api/internal/middleware"
	"go_api/internal/util"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// UserProfileHandler returns the user profile
func (h *UserHandler) UserProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get claims from context
		claims, ok := ctx.Value(middleware.UserClaimsKey).(*util.Claims)
		if !ok {
			util.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		// Get user profile from user service
		user, fromCache, err := h.service.GetUserProfile(ctx, claims.UserID)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound):
				util.ResponseWithError(w, http.StatusNotFound, "User not found", err.Error())
			case errors.Is(err, service.ErrCacheOperation):
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to cache user", err.Error())
			default:
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		message := "User profile (from database)"
		if fromCache {
			message = "User profile (from cache)"
		}
		util.ResponseWithSuccess(w, http.StatusOK, message, user)
	}
}

// CreateUserHandler creates a new user
func (h *UserHandler) CreateUserHandler() http.HandlerFunc {
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

		// Create user in user service
		user, err := h.service.CreateUser(ctx, req)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrPasswordHashing):
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to hash password", err.Error())
			case errors.Is(err, service.ErrUserCreation):
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to create user", err.Error())
			default:
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		util.ResponseWithSuccess(w, http.StatusCreated, "User created successfully", user)
	}
}

// LoginUserHandler logs in a user
func (h *UserHandler) LoginUserHandler() http.HandlerFunc {
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

		// Login user in user service
		token, err := h.service.LoginUser(ctx, req)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrUserNotFound):
				util.ResponseWithError(w, http.StatusNotFound, "User not found", err.Error())
			case errors.Is(err, service.ErrInvalidPassword):
				util.ResponseWithError(w, http.StatusUnauthorized, "Invalid password", "")
			case errors.Is(err, service.ErrTokenGeneration):
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
			default:
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "Login successful", map[string]string{
			"token": token,
		})
	}
}

// LogoutUserHandler logs out a user
func (h *UserHandler) LogoutUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get claims from context
		claims, ok := ctx.Value(middleware.UserClaimsKey).(*util.Claims)
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

		token, err := util.ExtractTokenFromHeader(authHeader)
		if err != nil {
			util.ResponseWithError(w, http.StatusUnauthorized, "Missing authorization token", err.Error())
			return
		}

		// Logout user in user service
		err = h.service.LogoutUser(ctx, claims.UserID, token, claims.ExpiresAt)
		if err != nil {
			switch {
			case errors.Is(err, service.ErrTokenBlacklist):
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to blacklist token", err.Error())
			case errors.Is(err, service.ErrSessionCleanup):
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to clean user session", err.Error())
			default:
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "Logout successful", nil)
	}
}

// ListAllUsersHandler lists all users (with optional pagination)
func (h *UserHandler) ListAllUsersHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// List all users from user service
		users, err := h.service.ListAllUsers(ctx)
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to list users", err.Error())
			return
		}
		util.ResponseWithSuccess(w, http.StatusOK, "List of all users", users)
	}
}


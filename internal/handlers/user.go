package handlers

import (
	"encoding/json"
	"net/http"

	"go_api/internal/auth"
	"go_api/internal/config"
	"go_api/internal/dto"
	"go_api/internal/models"
	"go_api/internal/utils"
)

func (h *Handler) UserProfileHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Get claims from context
		claims, ok := r.Context().Value(auth.UserClaimsKey).(*utils.Claims)
		if !ok {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		// Get user from database
		user := &models.User{}
		if err := h.DB.WithContext(ctx).Where("id = ?", claims.UserID).First(user).Error; err != nil {
			utils.ResponseWithError(w, http.StatusNotFound, "User not found", err.Error())
			return
		}

		utils.ResponseWithSuccess(w, http.StatusOK, "User profile", user)
	}
}

func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Validate request
		if err := utils.Validate(req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Create user
		hashedPassword, err := utils.HashPassword(req.Password)
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to hash password", err.Error())
			return
		}

		user := &models.User{
			Username: req.Username,
			Email:    req.Email,
			Password: hashedPassword,
		}

		if err := h.DB.WithContext(ctx).Create(user).Error; err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to create user", err.Error())
			return
		}

		utils.ResponseWithSuccess(w, http.StatusCreated, "User created successfully", user)
	}
}

func (h *Handler) LoginUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.LoginUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Validate request
		if err := utils.Validate(req); err != nil {
			utils.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Fetch user from DB
		user := &models.User{}
		if err := h.DB.WithContext(ctx).Where("email = ?", req.Email).First(user).Error; err != nil {
			utils.ResponseWithError(w, http.StatusNotFound, "User not found", err.Error())
			return
		}

		// Compare password
		if !utils.ComparePassword(req.Password, user.Password) {
			utils.ResponseWithError(w, http.StatusUnauthorized, "Invalid password", "")
			return
		}

		// Generate token
		token, err := utils.GenerateToken(user.ID, user.Username, []byte(config.GlobalConfig.JWTSecretKey))
		if err != nil {
			utils.ResponseWithError(w, http.StatusInternalServerError, "Failed to generate token", err.Error())
			return
		}

		// Return token
		utils.ResponseWithSuccess(w, http.StatusOK, "Login successful", map[string]string{
			"token": token,
		})
	}
}

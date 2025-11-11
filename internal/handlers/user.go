package handlers

import (
	"encoding/json"
	"go_api/internal/dto"
	"go_api/internal/models"
	"go_api/internal/utils"
	"net/http"
)

func (h *Handler) CreateUserHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var req dto.CreateUserRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
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
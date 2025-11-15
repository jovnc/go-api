package handler

import (
	"encoding/json"
	"net/http"

	"go_api/internal/app/dto"
	"go_api/internal/app/model"
	"go_api/internal/middleware"
	"go_api/internal/util"
)

// CreateBlog creates a new blog
func (h *Handler) CreateBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// Decode and validate request body
		var req dto.CreateBlogRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}
		if err := util.Validate(req); err != nil {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request body", err.Error())
			return
		}

		// Get claims from context
		claims, ok := ctx.Value(middleware.UserClaimsKey).(*util.Claims)
		if !ok {
			util.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		blog := &model.Blog{
			Title:   req.Title,
			Content: req.Content,
			UserID:  claims.UserID,
		}
		if err := h.DB.WithContext(ctx).Create(blog).Error; err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to create blog", err.Error())
			return
		}

		util.ResponseWithSuccess(w, http.StatusCreated, "Blog created successfully", blog)
	}
}

// GetBlogHandler gets a blog by its ID
func (h *Handler) GetBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request", "Blog ID is required")
			return
		}

		blog := &model.Blog{}
		if err := h.DB.WithContext(ctx).First(blog, id).Error; err != nil {
			util.ResponseWithError(w, http.StatusNotFound, "Blog not found", err.Error())
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "Blog retrieved successfully", blog)
	}
}

// DeleteBlogHandler deletes a blog by its ID
func (h *Handler) DeleteBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request", "Blog ID is required")
			return
		}

		claims, ok := ctx.Value(middleware.UserClaimsKey).(*util.Claims)
		if !ok {
			util.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		result := h.DB.WithContext(ctx).Where("id = ? AND user_id = ?", id, claims.UserID).Delete(&model.Blog{})
		if result.Error != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to delete blog", result.Error.Error())
			return
		}

		if result.RowsAffected == 0 {
			util.ResponseWithError(w, http.StatusNotFound, "Blog not found", "Blog not found or you don't have permission to delete it")
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "Blog deleted successfully", nil)
	}
}

// ListBlogsHandler lists all blogs
func (h *Handler) ListBlogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		blogs := []model.Blog{}
		if err := h.DB.WithContext(ctx).Find(&blogs).Error; err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to list blogs", err.Error())
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "List of all blogs", blogs)
	}
}

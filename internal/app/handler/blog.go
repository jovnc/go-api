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

		// Create blog
		blog := &model.Blog{
			Title:   req.Title,
			Content: req.Content,
			UserID:  claims.UserID,
		}
		if err := h.BlogRepository.CreateBlog(ctx, blog); err != nil {
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

		blog, err := h.BlogRepository.GetBlog(ctx, id)
		if err != nil {
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

		// Get claims from context
		claims, ok := ctx.Value(middleware.UserClaimsKey).(*util.Claims)
		if !ok {
			util.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		// Delete blog
		result := h.BlogRepository.DeleteBlog(ctx, id, claims.UserID)
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

		blogs, err := h.BlogRepository.ListBlogs(ctx)
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to list blogs", err.Error())
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "List of all blogs", blogs)
	}
}

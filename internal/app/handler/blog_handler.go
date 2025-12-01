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

type BlogHandler struct {
	service *service.BlogService
}

func NewBlogHandler(service *service.BlogService) *BlogHandler {
	return &BlogHandler{
		service: service,
	}
}

// CreateBlogHandler creates a new blog
func (h *BlogHandler) CreateBlogHandler() http.HandlerFunc {
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
		claims, ok := ctx.Value(middleware.UserClaimsKey).(*util.UserClaims)
		if !ok {
			util.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		// Create blog in blog service
		blog, err := h.service.CreateBlog(ctx, req, claims.UserID)
		if err != nil {
			if errors.Is(err, service.ErrBlogCreation) {
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to create blog", err.Error())
			} else {
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		util.ResponseWithSuccess(w, http.StatusCreated, "Blog created successfully", blog)
	}
}

// GetBlogHandler gets a blog by its ID
func (h *BlogHandler) GetBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request", "Blog ID is required")
			return
		}

		// Get blog from blog service
		blog, err := h.service.GetBlog(ctx, id)
		if err != nil {
			if errors.Is(err, service.ErrBlogNotFound) {
				util.ResponseWithError(w, http.StatusNotFound, "Blog not found", err.Error())
			} else {
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "Blog retrieved successfully", blog)
	}
}

// DeleteBlogHandler deletes a blog by its ID
func (h *BlogHandler) DeleteBlogHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		id := r.PathValue("id")
		if id == "" {
			util.ResponseWithError(w, http.StatusBadRequest, "Invalid request", "Blog ID is required")
			return
		}

		// Get claims from context
		claims, ok := ctx.Value(middleware.UserClaimsKey).(*util.UserClaims)
		if !ok {
			util.ResponseWithError(w, http.StatusUnauthorized, "Unauthorized", "Unauthorized")
			return
		}

		// Delete blog from blog service
		err := h.service.DeleteBlog(ctx, id, claims.UserID)
		if err != nil {
			if errors.Is(err, service.ErrBlogDeletion) {
				util.ResponseWithError(w, http.StatusInternalServerError, "Failed to delete blog", err.Error())
			} else {
				util.ResponseWithError(w, http.StatusInternalServerError, "Internal server error", err.Error())
			}
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "Blog deleted successfully", nil)
	}
}

// ListBlogsHandler lists all blogs
func (h *BlogHandler) ListBlogsHandler() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		// List all blogs from blog service
		blogs, err := h.service.ListBlogs(ctx)
		if err != nil {
			util.ResponseWithError(w, http.StatusInternalServerError, "Failed to list blogs", err.Error())
			return
		}

		util.ResponseWithSuccess(w, http.StatusOK, "List of all blogs", blogs)
	}
}

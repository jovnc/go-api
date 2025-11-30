package service

import (
	"context"
	"errors"

	"go_api/internal/app/dto"
	"go_api/internal/app/model"
	"go_api/internal/app/repository"

	"gorm.io/gorm"
)

var (
	ErrBlogNotFound   = errors.New("blog not found")
	ErrBlogCreation   = errors.New("failed to create blog")
	ErrBlogDeletion   = errors.New("failed to delete blog")
	ErrBlogListFailed = errors.New("failed to list blogs")
)

type BlogService struct {
	repo *repository.BlogRepository
}

func NewBlogService(db *gorm.DB) *BlogService {
	return &BlogService{
		repo: repository.NewBlogRepository(db),
	}
}

// CreateBlog creates a new blog post
func (s *BlogService) CreateBlog(ctx context.Context, req dto.CreateBlogRequest, userID uint) (*model.Blog, error) {
	blog := &model.Blog{
		Title:   req.Title,
		Content: req.Content,
		UserID:  userID,
	}

	if err := s.repo.CreateBlog(ctx, blog); err != nil {
		return nil, ErrBlogCreation
	}

	return blog, nil
}

// GetBlog retrieves a blog by its ID
func (s *BlogService) GetBlog(ctx context.Context, id string) (*model.Blog, error) {
	blog, err := s.repo.GetBlog(ctx, id)
	if err != nil {
		return nil, ErrBlogNotFound
	}
	return blog, nil
}

// DeleteBlog deletes a blog by its ID (only if owned by user)
func (s *BlogService) DeleteBlog(ctx context.Context, id string, userID uint) error {
	if err := s.repo.DeleteBlog(ctx, id, userID); err != nil {
		return ErrBlogDeletion
	}
	return nil
}

// ListBlogs retrieves all blogs
func (s *BlogService) ListBlogs(ctx context.Context) ([]model.Blog, error) {
	blogs, err := s.repo.ListBlogs(ctx)
	if err != nil {
		return nil, ErrBlogListFailed
	}
	return blogs, nil
}

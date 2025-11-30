package repository

import (
	"context"
	"errors"
	"go_api/internal/app/model"

	"gorm.io/gorm"
)

type BlogRepository struct {
	db *gorm.DB
}

func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

func (r *BlogRepository) CreateBlog(ctx context.Context, blog *model.Blog) error {
	return r.db.WithContext(ctx).Create(blog).Error
}

func (r *BlogRepository) GetBlog(ctx context.Context, id string) (*model.Blog, error) {
	var blog model.Blog
	if err := r.db.WithContext(ctx).First(&blog, id).Error; err != nil {
		return nil, err
	}
	return &blog, nil
}

func (r *BlogRepository) DeleteBlog(ctx context.Context, id string, userID uint) error {
	result := r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Blog{})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return errors.New("blog not found")
	}
	return nil
}

func (r *BlogRepository) ListBlogs(ctx context.Context) ([]model.Blog, error) {
	var blogs []model.Blog
	if err := r.db.WithContext(ctx).Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}



package repository

import (
	"context"
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

func (r *BlogRepository) DeleteBlog(ctx context.Context, id string, userID uint) *gorm.DB {
	return r.db.WithContext(ctx).Where("id = ? AND user_id = ?", id, userID).Delete(&model.Blog{})
}

func (r *BlogRepository) ListBlogs(ctx context.Context) ([]model.Blog, error) {
	var blogs []model.Blog
	if err := r.db.WithContext(ctx).Find(&blogs).Error; err != nil {
		return nil, err
	}
	return blogs, nil
}



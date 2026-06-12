package repository

import (
	"context"
	"fmt"

	"go-cms/internal/model"

	"gorm.io/gorm"
)

// BlogRepository handles all blog database queries.
type BlogRepository struct {
	db *gorm.DB
}

// NewBlogRepository creates a BlogRepository with the given GORM instance.
func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

// FindAll retrieves all blog posts.
func (r *BlogRepository) FindAll(ctx context.Context) ([]model.Blog, error) {
	var blogs []model.Blog
	err := r.db.WithContext(ctx).Order("created_at desc").Find(&blogs).Error
	if err != nil {
		return nil, fmt.Errorf("blog repo: find all: %w", err)
	}
	return blogs, nil
}

// FindAllPublished retrieves all published blog posts.
func (r *BlogRepository) FindAllPublished(ctx context.Context) ([]model.Blog, error) {
	var blogs []model.Blog
	err := r.db.WithContext(ctx).Where("status = ?", "published").Order("published_at desc").Find(&blogs).Error
	if err != nil {
		return nil, fmt.Errorf("blog repo: find all published: %w", err)
	}
	return blogs, nil
}

// FindBySlug retrieves a blog post by its unique slug.
// Returns (nil, nil) if not found.
func (r *BlogRepository) FindBySlug(ctx context.Context, slug string) (*model.Blog, error) {
	var blog model.Blog
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&blog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("blog repo: find by slug: %w", err)
	}
	return &blog, nil
}

// FindByID retrieves a blog post by its ID.
// Returns (nil, nil) if not found.
func (r *BlogRepository) FindByID(ctx context.Context, id string) (*model.Blog, error) {
	var blog model.Blog
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&blog).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("blog repo: find by id: %w", err)
	}
	return &blog, nil
}

// Create inserts a new blog post.
func (r *BlogRepository) Create(ctx context.Context, blog *model.Blog) error {
	if err := r.db.WithContext(ctx).Create(blog).Error; err != nil {
		return fmt.Errorf("blog repo: create: %w", err)
	}
	return nil
}

// Update updates an existing blog post.
func (r *BlogRepository) Update(ctx context.Context, blog *model.Blog) error {
	if err := r.db.WithContext(ctx).Save(blog).Error; err != nil {
		return fmt.Errorf("blog repo: update: %w", err)
	}
	return nil
}

// Delete removes a blog post by its slug.
func (r *BlogRepository) Delete(ctx context.Context, slug string) error {
	if err := r.db.WithContext(ctx).Where("slug = ?", slug).Delete(&model.Blog{}).Error; err != nil {
		return fmt.Errorf("blog repo: delete by slug: %w", err)
	}
	return nil
}

// Count returns the total number of blog posts.
func (r *BlogRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Blog{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("blog repo: count: %w", err)
	}
	return count, nil
}

// Recent retrieves the most recent published blog posts up to the specified limit.
func (r *BlogRepository) Recent(ctx context.Context, limit int) ([]model.Blog, error) {
	var blogs []model.Blog
	err := r.db.WithContext(ctx).Where("status = ?", "published").Order("published_at desc").Limit(limit).Find(&blogs).Error
	if err != nil {
		return nil, fmt.Errorf("blog repo: recent: %w", err)
	}
	return blogs, nil
}


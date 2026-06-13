package repository

import (
	"context"
	"fmt"
	"time"

	"go-cms/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// BlogRepository handles all blog database queries using Raw SQL.
type BlogRepository struct {
	db *gorm.DB
}

// NewBlogRepository creates a BlogRepository with the given GORM instance.
func NewBlogRepository(db *gorm.DB) *BlogRepository {
	return &BlogRepository{db: db}
}

// FindAll retrieves all blog posts.
func (r *BlogRepository) FindAll(ctx context.Context) ([]dto.BlogResponse, error) {
	var blogs []dto.BlogResponse
	query := `
		SELECT id, title, slug, category, content, excerpt, status, published_at, created_at, updated_at 
		FROM blogs 
		ORDER BY created_at DESC
	`
	err := r.db.WithContext(ctx).Raw(query).Scan(&blogs).Error
	if err != nil {
		return nil, fmt.Errorf("blog repo: find all: %w", err)
	}
	return blogs, nil
}

// FindAllPublished retrieves all published blog posts.
func (r *BlogRepository) FindAllPublished(ctx context.Context) ([]dto.BlogResponse, error) {
	var blogs []dto.BlogResponse
	query := `
		SELECT id, title, slug, category, content, excerpt, status, published_at, created_at, updated_at 
		FROM blogs 
		WHERE status = 'published' 
		ORDER BY published_at DESC
	`
	err := r.db.WithContext(ctx).Raw(query).Scan(&blogs).Error
	if err != nil {
		return nil, fmt.Errorf("blog repo: find all published: %w", err)
	}
	return blogs, nil
}

// FindBySlug retrieves a blog post by its unique slug.
func (r *BlogRepository) FindBySlug(ctx context.Context, slug string) (dto.BlogResponse, error) {
	var resp dto.BlogResponse
	query := `
		SELECT id, title, slug, category, content, excerpt, status, published_at, created_at, updated_at 
		FROM blogs 
		WHERE slug = ?
	`
	err := r.db.WithContext(ctx).Raw(query, slug).Scan(&resp).Error
	if err != nil {
		return dto.BlogResponse{}, fmt.Errorf("blog repo: find by slug: %w", err)
	}

	if resp.ID == "" {
		return dto.BlogResponse{}, gorm.ErrRecordNotFound
	}
	return resp, nil
}

// FindByID retrieves a blog post by its ID.
func (r *BlogRepository) FindByID(ctx context.Context, id string) (dto.BlogResponse, error) {
	var resp dto.BlogResponse
	query := `
		SELECT id, title, slug, category, content, excerpt, status, published_at, created_at, updated_at 
		FROM blogs 
		WHERE id = ?
	`
	err := r.db.WithContext(ctx).Raw(query, id).Scan(&resp).Error
	if err != nil {
		return dto.BlogResponse{}, fmt.Errorf("blog repo: find by id: %w", err)
	}

	if resp.ID == "" {
		return dto.BlogResponse{}, gorm.ErrRecordNotFound
	}
	return resp, nil
}

// Create inserts a new blog post.
func (r *BlogRepository) Create(ctx context.Context, req dto.CreateBlogRequest) (dto.BlogResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return dto.BlogResponse{}, fmt.Errorf("blog repo: create: generate uuid: %w", err)
	}
	now := time.Now()

	var publishedAt *time.Time
	if req.Status == "published" {
		if req.PublishedAt != nil {
			publishedAt = req.PublishedAt
		} else {
			publishedAt = &now
		}
	}

	query := `
		INSERT INTO blogs (id, title, slug, category, content, excerpt, status, published_at, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	err = r.db.WithContext(ctx).Exec(query, id, req.Title, req.Slug, req.Category, req.Content, req.Excerpt, req.Status, publishedAt, now, now).Error
	if err != nil {
		return dto.BlogResponse{}, fmt.Errorf("blog repo: create: %w", err)
	}

	return dto.BlogResponse{
		ID:          id.String(),
		Title:       req.Title,
		Slug:        req.Slug,
		Category:    req.Category,
		Content:     req.Content,
		Excerpt:     req.Excerpt,
		Status:      req.Status,
		PublishedAt: publishedAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update updates an existing blog post.
func (r *BlogRepository) Update(ctx context.Context, id string, req dto.UpdateBlogRequest) (dto.BlogResponse, error) {
	now := time.Now()

	var publishedAt *time.Time
	if req.Status == "published" {
		if req.PublishedAt != nil {
			publishedAt = req.PublishedAt
		} else {
			publishedAt = &now
		}
	}

	query := `
		UPDATE blogs 
		SET title = ?, slug = ?, category = ?, content = ?, excerpt = ?, status = ?, published_at = ?, updated_at = ? 
		WHERE id = ?
	`
	result := r.db.WithContext(ctx).Exec(query, req.Title, req.Slug, req.Category, req.Content, req.Excerpt, req.Status, publishedAt, now, id)
	if result.Error != nil {
		return dto.BlogResponse{}, fmt.Errorf("blog repo: update: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return dto.BlogResponse{}, gorm.ErrRecordNotFound
	}

	return r.FindByID(ctx, id)
}

// Delete removes a blog post by its slug.
func (r *BlogRepository) Delete(ctx context.Context, slug string) error {
	query := `DELETE FROM blogs WHERE slug = ?`
	result := r.db.WithContext(ctx).Exec(query, slug)
	if result.Error != nil {
		return fmt.Errorf("blog repo: delete by slug: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Count returns the total number of blog posts.
func (r *BlogRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM blogs`
	err := r.db.WithContext(ctx).Raw(query).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("blog repo: count: %w", err)
	}
	return count, nil
}

// Recent retrieves the most recent published blog posts up to the specified limit.
func (r *BlogRepository) Recent(ctx context.Context, limit int) ([]dto.BlogResponse, error) {
	var blogs []dto.BlogResponse
	query := `
		SELECT id, title, slug, category, content, excerpt, status, published_at, created_at, updated_at 
		FROM blogs 
		WHERE status = 'published' 
		ORDER BY published_at DESC 
		LIMIT ?
	`
	err := r.db.WithContext(ctx).Raw(query, limit).Scan(&blogs).Error
	if err != nil {
		return nil, fmt.Errorf("blog repo: recent: %w", err)
	}
	return blogs, nil
}

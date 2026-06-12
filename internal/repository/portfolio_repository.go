package repository

import (
	"context"
	"fmt"

	"go-cms/internal/model"

	"gorm.io/gorm"
)

// PortfolioRepository handles all portfolio database queries.
type PortfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository creates a PortfolioRepository with the given GORM instance.
func NewPortfolioRepository(db *gorm.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

// FindAll retrieves all portfolio items ordered by sort_order.
func (r *PortfolioRepository) FindAll(ctx context.Context) ([]model.Portfolio, error) {
	var portfolios []model.Portfolio
	err := r.db.WithContext(ctx).Order("sort_order asc, created_at desc").Find(&portfolios).Error
	if err != nil {
		return nil, fmt.Errorf("portfolio repo: find all: %w", err)
	}
	return portfolios, nil
}

// FindByID retrieves a portfolio item by its ID.
// Returns (nil, nil) if not found.
func (r *PortfolioRepository) FindByID(ctx context.Context, id string) (*model.Portfolio, error) {
	var portfolio model.Portfolio
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&portfolio).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("portfolio repo: find by id: %w", err)
	}
	return &portfolio, nil
}

// FindByTitle retrieves a portfolio item by its exact title.
// Returns (nil, nil) if not found.
func (r *PortfolioRepository) FindByTitle(ctx context.Context, title string) (*model.Portfolio, error) {
	var portfolio model.Portfolio
	err := r.db.WithContext(ctx).Where("title = ?", title).First(&portfolio).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("portfolio repo: find by title: %w", err)
	}
	return &portfolio, nil
}

// Create inserts a new portfolio item.
func (r *PortfolioRepository) Create(ctx context.Context, portfolio *model.Portfolio) error {
	if err := r.db.WithContext(ctx).Create(portfolio).Error; err != nil {
		return fmt.Errorf("portfolio repo: create: %w", err)
	}
	return nil
}

// Update updates an existing portfolio item.
func (r *PortfolioRepository) Update(ctx context.Context, portfolio *model.Portfolio) error {
	if err := r.db.WithContext(ctx).Save(portfolio).Error; err != nil {
		return fmt.Errorf("portfolio repo: update: %w", err)
	}
	return nil
}

// Delete removes a portfolio item by its title.
func (r *PortfolioRepository) Delete(ctx context.Context, title string) error {
	if err := r.db.WithContext(ctx).Where("title = ?", title).Delete(&model.Portfolio{}).Error; err != nil {
		return fmt.Errorf("portfolio repo: delete by title: %w", err)
	}
	return nil
}

// Count returns the total number of portfolio items.
func (r *PortfolioRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	if err := r.db.WithContext(ctx).Model(&model.Portfolio{}).Count(&count).Error; err != nil {
		return 0, fmt.Errorf("portfolio repo: count: %w", err)
	}
	return count, nil
}


package repository

import (
	"context"
	"fmt"
	"time"

	"go-cms/internal/dto"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// PortfolioRepository handles all portfolio database queries using Raw SQL.
type PortfolioRepository struct {
	db *gorm.DB
}

// NewPortfolioRepository creates a PortfolioRepository with the given GORM instance.
func NewPortfolioRepository(db *gorm.DB) *PortfolioRepository {
	return &PortfolioRepository{db: db}
}

// FindAll retrieves all portfolio items ordered by sort_order.
func (r *PortfolioRepository) FindAll(ctx context.Context) ([]dto.PortfolioResponse, error) {
	var results []struct {
		ID            uuid.UUID
		Title         string
		Icon          string
		Description   string
		TechStack     pq.StringArray
		ProjectURL    string
		RepositoryURL string
		SortOrder     int
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	query := `
		SELECT id, title, icon, description, tech_stack, project_url, repository_url, sort_order, created_at, updated_at 
		FROM portfolios 
		ORDER BY sort_order ASC, created_at DESC
	`
	err := r.db.WithContext(ctx).Raw(query).Scan(&results).Error
	if err != nil {
		return nil, fmt.Errorf("portfolio repo: find all: %w", err)
	}

	resp := make([]dto.PortfolioResponse, len(results))
	for i, p := range results {
		resp[i] = dto.PortfolioResponse{
			ID:            p.ID.String(),
			Title:         p.Title,
			Icon:          p.Icon,
			Description:   p.Description,
			TechStack:     []string(p.TechStack),
			ProjectURL:    p.ProjectURL,
			RepositoryURL: p.RepositoryURL,
			SortOrder:     p.SortOrder,
			CreatedAt:     p.CreatedAt,
			UpdatedAt:     p.UpdatedAt,
		}
	}
	return resp, nil
}

// FindByID retrieves a portfolio item by its ID.
func (r *PortfolioRepository) FindByID(ctx context.Context, id string) (dto.PortfolioResponse, error) {
	var p struct {
		ID            uuid.UUID
		Title         string
		Icon          string
		Description   string
		TechStack     pq.StringArray
		ProjectURL    string
		RepositoryURL string
		SortOrder     int
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	query := `
		SELECT id, title, icon, description, tech_stack, project_url, repository_url, sort_order, created_at, updated_at 
		FROM portfolios 
		WHERE id = ?
	`
	err := r.db.WithContext(ctx).Raw(query, id).Scan(&p).Error
	if err != nil {
		return dto.PortfolioResponse{}, fmt.Errorf("portfolio repo: find by id: %w", err)
	}

	if p.ID == uuid.Nil {
		return dto.PortfolioResponse{}, gorm.ErrRecordNotFound
	}

	return dto.PortfolioResponse{
		ID:            p.ID.String(),
		Title:         p.Title,
		Icon:          p.Icon,
		Description:   p.Description,
		TechStack:     []string(p.TechStack),
		ProjectURL:    p.ProjectURL,
		RepositoryURL: p.RepositoryURL,
		SortOrder:     p.SortOrder,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, nil
}

// FindByTitle retrieves a portfolio item by its exact title.
func (r *PortfolioRepository) FindByTitle(ctx context.Context, title string) (dto.PortfolioResponse, error) {
	var p struct {
		ID            uuid.UUID
		Title         string
		Icon          string
		Description   string
		TechStack     pq.StringArray
		ProjectURL    string
		RepositoryURL string
		SortOrder     int
		CreatedAt     time.Time
		UpdatedAt     time.Time
	}

	query := `
		SELECT id, title, icon, description, tech_stack, project_url, repository_url, sort_order, created_at, updated_at 
		FROM portfolios 
		WHERE title = ?
	`
	err := r.db.WithContext(ctx).Raw(query, title).Scan(&p).Error
	if err != nil {
		return dto.PortfolioResponse{}, fmt.Errorf("portfolio repo: find by title: %w", err)
	}

	if p.ID == uuid.Nil {
		return dto.PortfolioResponse{}, gorm.ErrRecordNotFound
	}

	return dto.PortfolioResponse{
		ID:            p.ID.String(),
		Title:         p.Title,
		Icon:          p.Icon,
		Description:   p.Description,
		TechStack:     []string(p.TechStack),
		ProjectURL:    p.ProjectURL,
		RepositoryURL: p.RepositoryURL,
		SortOrder:     p.SortOrder,
		CreatedAt:     p.CreatedAt,
		UpdatedAt:     p.UpdatedAt,
	}, nil
}

// Create inserts a new portfolio item.
func (r *PortfolioRepository) Create(ctx context.Context, req dto.CreatePortfolioRequest) (dto.PortfolioResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return dto.PortfolioResponse{}, fmt.Errorf("portfolio repo: create: generate uuid: %w", err)
	}
	now := time.Now()

	query := `
		INSERT INTO portfolios (id, title, icon, description, tech_stack, project_url, repository_url, sort_order, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	err = r.db.WithContext(ctx).Exec(query, id, req.Title, req.Icon, req.Description, pq.StringArray(req.TechStack), req.ProjectURL, req.RepositoryURL, req.SortOrder, now, now).Error
	if err != nil {
		return dto.PortfolioResponse{}, fmt.Errorf("portfolio repo: create: %w", err)
	}

	return dto.PortfolioResponse{
		ID:            id.String(),
		Title:         req.Title,
		Icon:          req.Icon,
		Description:   req.Description,
		TechStack:     req.TechStack,
		ProjectURL:    req.ProjectURL,
		RepositoryURL: req.RepositoryURL,
		SortOrder:     req.SortOrder,
		CreatedAt:     now,
		UpdatedAt:     now,
	}, nil
}

// Update updates an existing portfolio item.
func (r *PortfolioRepository) Update(ctx context.Context, id string, req dto.UpdatePortfolioRequest) (dto.PortfolioResponse, error) {
	now := time.Now()

	query := `
		UPDATE portfolios 
		SET title = ?, icon = ?, description = ?, tech_stack = ?, project_url = ?, repository_url = ?, sort_order = ?, updated_at = ? 
		WHERE id = ?
	`
	result := r.db.WithContext(ctx).Exec(query, req.Title, req.Icon, req.Description, pq.StringArray(req.TechStack), req.ProjectURL, req.RepositoryURL, req.SortOrder, now, id)
	if result.Error != nil {
		return dto.PortfolioResponse{}, fmt.Errorf("portfolio repo: update: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return dto.PortfolioResponse{}, gorm.ErrRecordNotFound
	}

	return r.FindByID(ctx, id)
}

// Delete removes a portfolio item by its title.
func (r *PortfolioRepository) Delete(ctx context.Context, title string) error {
	query := `DELETE FROM portfolios WHERE title = ?`
	result := r.db.WithContext(ctx).Exec(query, title)
	if result.Error != nil {
		return fmt.Errorf("portfolio repo: delete by title: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound
	}
	return nil
}

// Count returns the total number of portfolio items.
func (r *PortfolioRepository) Count(ctx context.Context) (int64, error) {
	var count int64
	query := `SELECT COUNT(*) FROM portfolios`
	err := r.db.WithContext(ctx).Raw(query).Scan(&count).Error
	if err != nil {
		return 0, fmt.Errorf("portfolio repo: count: %w", err)
	}
	return count, nil
}

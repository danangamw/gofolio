package dto

import (
	"time"
)

type CreatePortfolioRequest struct {
	Title         string   `json:"title" validate:"required,max=255"`
	Icon          string   `json:"icon" validate:"required,max=50"`
	Description   string   `json:"description" validate:"required"`
	TechStack     []string `json:"tech_stack" validate:"required"`
	ProjectURL    string   `json:"project_url" validate:"max=500"`
	RepositoryURL string   `json:"repository_url" validate:"max=500"`
	SortOrder     int      `json:"sort_order"`
}

type UpdatePortfolioRequest struct {
	Title         string   `json:"title" validate:"required,max=255"`
	Icon          string   `json:"icon" validate:"required,max=50"`
	Description   string   `json:"description" validate:"required"`
	TechStack     []string `json:"tech_stack" validate:"required"`
	ProjectURL    string   `json:"project_url" validate:"max=500"`
	RepositoryURL string   `json:"repository_url" validate:"max=500"`
	SortOrder     int      `json:"sort_order"`
}

type PortfolioResponse struct {
	ID            string    `json:"id"`
	Title         string    `json:"title"`
	Icon          string    `json:"icon"`
	Description   string    `json:"description"`
	TechStack     []string  `json:"tech_stack"`
	ProjectURL    string    `json:"project_url"`
	RepositoryURL string    `json:"repository_url"`
	SortOrder     int       `json:"sort_order"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

package repository

import (
	"context"
	"fmt"
	"time"

	"go-cms/internal/dto"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// UserRepository handles all user-related database queries using Raw SQL.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a UserRepository with the given GORM instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUsername retrieves a user by their username.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (dto.UserResponse, error) {
	var resp dto.UserResponse
	query := `
		SELECT id, username, password_hash, created_at, updated_at 
		FROM users 
		WHERE username = ?
	`
	err := r.db.WithContext(ctx).Raw(query, username).Scan(&resp).Error
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("user repo: find by username: %w", err)
	}

	if resp.ID == "" {
		return dto.UserResponse{}, gorm.ErrRecordNotFound
	}

	return resp, nil
}

// FindByID retrieves a user by their UUID.
func (r *UserRepository) FindByID(ctx context.Context, id string) (dto.UserResponse, error) {
	var resp dto.UserResponse
	query := `
		SELECT id, username, password_hash, created_at, updated_at 
		FROM users 
		WHERE id = ?
	`
	err := r.db.WithContext(ctx).Raw(query, id).Scan(&resp).Error
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("user repo: find by id: %w", err)
	}

	if resp.ID == "" {
		return dto.UserResponse{}, gorm.ErrRecordNotFound
	}

	return resp, nil
}

// Create inserts a new user record.
func (r *UserRepository) Create(ctx context.Context, req dto.CreateUserRequest) (dto.UserResponse, error) {
	id, err := uuid.NewV7()
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("user repo: create: generate uuid: %w", err)
	}
	now := time.Now()

	query := `
		INSERT INTO users (id, username, password_hash, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?)
	`
	err = r.db.WithContext(ctx).Exec(query, id, req.Username, req.PasswordHash, now, now).Error
	if err != nil {
		return dto.UserResponse{}, fmt.Errorf("user repo: create: %w", err)
	}

	return dto.UserResponse{
		ID:           id.String(),
		Username:     req.Username,
		PasswordHash: req.PasswordHash,
		CreatedAt:    now,
		UpdatedAt:    now,
	}, nil
}

package repository

import (
	"context"
	"fmt"

	"go-cms/internal/model"

	"gorm.io/gorm"
)

// UserRepository handles all user-related database queries.
type UserRepository struct {
	db *gorm.DB
}

// NewUserRepository creates a UserRepository with the given GORM instance.
func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

// FindByUsername retrieves a user by their username.
// Returns (nil, nil) if no user is found.
func (r *UserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("username = ?", username).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("user repo: find by username: %w", err)
	}
	return &user, nil
}

// FindByID retrieves a user by their UUID.
// Returns (nil, nil) if not found.
func (r *UserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	var user model.User
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, fmt.Errorf("user repo: find by id: %w", err)
	}
	return &user, nil
}

// Create inserts a new user record and returns the created user.
func (r *UserRepository) Create(ctx context.Context, user *model.User) error {
	if err := r.db.WithContext(ctx).Create(user).Error; err != nil {
		return fmt.Errorf("user repo: create: %w", err)
	}
	return nil
}

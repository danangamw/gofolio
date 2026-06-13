package dto

import (
	"time"
)

type CreateUserRequest struct {
	Username     string `json:"username" validate:"required,min=3,max=255"`
	PasswordHash string `json:"password_hash" validate:"required"`
}

type UserResponse struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	PasswordHash string    `json:"-"` // Excluded from JSON serialization
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

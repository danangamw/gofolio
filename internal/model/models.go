package model

import (
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"gorm.io/gorm"
)

// User represents the users table in database.
type User struct {
	ID           uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Username     string         `gorm:"type:varchar(100);uniqueIndex;not null"`
	PasswordHash string         `gorm:"type:text;not null"`
	CreatedAt    time.Time      `gorm:"autoCreateTime"`
	UpdatedAt    time.Time      `gorm:"autoUpdateTime"`
	Sessions     []Session      `gorm:"foreignKey:UserID;constraint:OnDelete:CASCADE"`
}

// Blog represents the blogs table in database.
type Blog struct {
	ID          uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title       string         `gorm:"type:varchar(255);not null"`
	Slug        string         `gorm:"type:varchar(255);uniqueIndex;not null"`
	Content     string         `gorm:"type:text;not null"`
	Excerpt     string         `gorm:"type:text"`
	Status      string         `gorm:"type:varchar(20);default:'draft';not null"` // 'draft' or 'published'
	PublishedAt *time.Time     `gorm:"default:null"`
	CreatedAt   time.Time      `gorm:"autoCreateTime"`
	UpdatedAt   time.Time      `gorm:"autoUpdateTime"`
}

// Portfolio represents the portfolios table in database.
type Portfolio struct {
	ID            uuid.UUID      `gorm:"type:uuid;primaryKey;default:gen_random_uuid()"`
	Title         string         `gorm:"type:varchar(255);not null"`
	Description   string         `gorm:"type:text"`
	ImageURL      string         `gorm:"type:varchar(500)"`
	TechStack     pq.StringArray `gorm:"type:text[]"` // Native Postgres TEXT[] array
	ProjectURL    string         `gorm:"type:varchar(500)"`
	RepositoryURL string         `gorm:"type:varchar(500)"`
	SortOrder     int            `gorm:"type:integer;default:0;not null"`
	CreatedAt     time.Time      `gorm:"autoCreateTime"`
	UpdatedAt     time.Time      `gorm:"autoUpdateTime"`
}

// Session represents the server-side sessions table in database (fallback if Redis is not used).
type Session struct {
	ID           string    `gorm:"type:varchar(128);primaryKey"`
	UserID       uuid.UUID `gorm:"type:uuid;not null"`
	ExpiresAt    time.Time `gorm:"not null"`
	LastActiveAt time.Time `gorm:"not null"`
	CreatedAt    time.Time `gorm:"autoCreateTime"`
}

// GORM BeforeCreate Hooks for safe UUID generation in Go
func (u *User) BeforeCreate(tx *gorm.DB) (err error) {
	if u.ID == uuid.Nil {
		u.ID = uuid.New()
	}
	return
}

func (b *Blog) BeforeCreate(tx *gorm.DB) (err error) {
	if b.ID == uuid.Nil {
		b.ID = uuid.New()
	}
	return
}

func (p *Portfolio) BeforeCreate(tx *gorm.DB) (err error) {
	if p.ID == uuid.Nil {
		p.ID = uuid.New()
	}
	return
}

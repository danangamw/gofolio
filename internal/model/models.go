package model

import (
	"bytes"
	"html/template"
	"time"

	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
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
	Category    string         `gorm:"type:varchar(100);default:'';not null"`
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
	Icon          string         `gorm:"type:varchar(50);default:'';not null"`
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

// Date returns a formatted string of the publish/creation date.
func (b Blog) Date() string {
	if b.PublishedAt != nil {
		return b.PublishedAt.Format("January 2, 2006")
	}
	return b.CreatedAt.Format("January 2, 2006")
}

// HTMLContent converts markdown content to safe HTML.
func (b Blog) HTMLContent() template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(b.Content), &buf); err != nil {
		// Even for simple fallback, sanitize the output.
		sanitizedRaw := bluemonday.UGCPolicy().SanitizeBytes([]byte(b.Content))
		return template.HTML(sanitizedRaw)
	}
	sanitized := bluemonday.UGCPolicy().SanitizeBytes(buf.Bytes())
	return template.HTML(sanitized)
}

// Author returns a static author name 'Danang' as defined by user context.
func (b Blog) Author() string {
	return "Danang"
}



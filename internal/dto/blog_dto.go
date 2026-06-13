package dto

import (
	"bytes"
	"html/template"
	"time"

	"github.com/microcosm-cc/bluemonday"
	"github.com/yuin/goldmark"
)

type CreateBlogRequest struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Slug        string     `json:"slug" validate:"required,max=255"`
	Category    string     `json:"category" validate:"required,max=100"`
	Content     string     `json:"content" validate:"required"`
	Excerpt     string     `json:"excerpt" validate:"required,max=500"`
	Status      string     `json:"status" validate:"required,oneof=draft published"`
	PublishedAt *time.Time `json:"published_at"`
}

type UpdateBlogRequest struct {
	Title       string     `json:"title" validate:"required,max=255"`
	Slug        string     `json:"slug" validate:"required,max=255"`
	Category    string     `json:"category" validate:"required,max=100"`
	Content     string     `json:"content" validate:"required"`
	Excerpt     string     `json:"excerpt" validate:"required,max=500"`
	Status      string     `json:"status" validate:"required,oneof=draft published"`
	PublishedAt *time.Time `json:"published_at"`
}

type BlogResponse struct {
	ID          string     `json:"id"`
	Title       string     `json:"title"`
	Slug        string     `json:"slug"`
	Category    string     `json:"category"`
	Content     string     `json:"content"`
	Excerpt     string     `json:"excerpt"`
	Status      string     `json:"status"`
	PublishedAt *time.Time `json:"published_at"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
}

// Date returns a formatted string of the publish/creation date.
func (b BlogResponse) Date() string {
	if b.PublishedAt != nil {
		return b.PublishedAt.Format("January 2, 2006")
	}
	return b.CreatedAt.Format("January 2, 2006")
}

// HTMLContent converts markdown content to safe HTML.
func (b BlogResponse) HTMLContent() template.HTML {
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
func (b BlogResponse) Author() string {
	return "Danang"
}

package admin

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

type AdminBlogHandler struct {
	tmpl Renderer
}

func NewAdminBlogHandler(tmpl Renderer) *AdminBlogHandler {
	return &AdminBlogHandler{tmpl: tmpl}
}

type blogPostItem struct {
	Title       string
	Category    string
	Date        string
	Excerpt     string
	Slug        string
	HTMLContent string
}

var mockBlogs = []blogPostItem{
	{
		Category:    "Go",
		Date:        "June 10, 2026",
		Title:       "Understanding Dependency Injection in Go Simply",
		Excerpt:     "How to neatly manage database and third-party dependencies in Go applications without external frameworks.",
		Slug:        "understanding-dependency-injection-in-go",
		HTMLContent: "<p>Dependency Injection (DI) is often considered a complex concept...</p>",
	},
	{
		Category:    "Observability",
		Date:        "June 8, 2026",
		Title:       "Integrating OpenTelemetry Tracing in GORM",
		Excerpt:     "A complete guide on recording SQL query performance directly to Grafana Tempo using the OTel GORM plugin.",
		Slug:        "integrating-opentelemetry-tracing-in-gorm",
		HTMLContent: "<p>When our application slows down, one of the main suspects is...</p>",
	},
	{
		Category:    "Database",
		Date:        "June 5, 2026",
		Title:       "Designing Database Schema Migrations with Atlas",
		Excerpt:     "Why declarative migrations with Atlas are safer and more efficient than traditional manual SQL scripts.",
		Slug:        "designing-database-schema-migrations-with-atlas",
		HTMLContent: "<p>Managing database schema changes (migrations) often poses a major challenge...</p>",
	},
}

func (h *AdminBlogHandler) List(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "Manage Articles",
		"ActiveMenu": "blog",
		"Blogs":      mockBlogs,
	}
	h.tmpl.Render(w, "blog_list_admin", data)
}

func (h *AdminBlogHandler) New(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "New Article",
		"ActiveMenu": "blog",
		"IsEdit":     false,
	}
	h.tmpl.Render(w, "blog_form", data)
}

func (h *AdminBlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Simple redirect back to list for static demo
	http.Redirect(w, r, "/admin/blogs", http.StatusSeeOther)
}

func (h *AdminBlogHandler) Edit(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	var found blogPostItem
	for _, b := range mockBlogs {
		if b.Slug == slug {
			found = b
			break
		}
	}
	if found.Slug == "" {
		found = mockBlogs[0] // Fallback
	}

	data := map[string]any{
		"Title":      "Edit Article",
		"ActiveMenu": "blog",
		"IsEdit":     true,
		"Blog":       found,
	}
	h.tmpl.Render(w, "blog_form", data)
}

func (h *AdminBlogHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Simple redirect back to list for static demo
	http.Redirect(w, r, "/admin/blogs", http.StatusSeeOther)
}

func (h *AdminBlogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	// Simple redirect back to list for static demo
	slug := chi.URLParam(r, "slug")
	var newBlogs []blogPostItem
	for _, b := range mockBlogs {
		if b.Slug != slug {
			newBlogs = append(newBlogs, b)
		}
	}
	if len(newBlogs) < len(mockBlogs) {
		mockBlogs = newBlogs
	}
	http.Redirect(w, r, "/admin/blogs", http.StatusSeeOther)
}

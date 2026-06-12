package admin

import (
	"net/http"
	"regexp"
	"strings"
	"time"

	"go-cms/internal/middleware"
	"go-cms/internal/model"
	"go-cms/internal/repository"

	"github.com/go-chi/chi/v5"
)

type AdminBlogHandler struct {
	tmpl Renderer
	repo *repository.BlogRepository
}

func NewAdminBlogHandler(tmpl Renderer, repo *repository.BlogRepository) *AdminBlogHandler {
	return &AdminBlogHandler{
		tmpl: tmpl,
		repo: repo,
	}
}

func (h *AdminBlogHandler) List(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Manage Articles",
		"ActiveMenu": "blog",
		"Blogs":      blogs,
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "blog_list_admin", data)
}

func (h *AdminBlogHandler) New(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "New Article",
		"ActiveMenu": "blog",
		"IsEdit":     false,
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "blog_form", data)
}

func (h *AdminBlogHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	category := r.FormValue("category")
	dateStr := r.FormValue("date")
	excerpt := r.FormValue("excerpt")
	content := r.FormValue("content")

	blog := &model.Blog{
		Title:       title,
		Slug:        slugify(title),
		Category:    category,
		Excerpt:     excerpt,
		Content:     content,
		Status:      "published", // default to published since there's no status dropdown
		PublishedAt: parseDateInput(dateStr),
	}

	if err := h.repo.Create(r.Context(), blog); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/blogs", http.StatusSeeOther)
}

func (h *AdminBlogHandler) Edit(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	blog, err := h.repo.FindBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if blog == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title":      "Edit Article",
		"ActiveMenu": "blog",
		"IsEdit":     true,
		"Blog":       blog,
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "blog_form", data)
}

func (h *AdminBlogHandler) Update(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	blog, err := h.repo.FindBySlug(r.Context(), slug)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if blog == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	category := r.FormValue("category")
	dateStr := r.FormValue("date")
	excerpt := r.FormValue("excerpt")
	content := r.FormValue("content")

	blog.Title = title
	blog.Slug = slugify(title)
	blog.Category = category
	blog.Excerpt = excerpt
	blog.Content = content
	blog.PublishedAt = parseDateInput(dateStr)

	if err := h.repo.Update(r.Context(), blog); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/blogs", http.StatusSeeOther)
}

func (h *AdminBlogHandler) Delete(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	if err := h.repo.Delete(r.Context(), slug); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/blogs", http.StatusSeeOther)
}

func slugify(title string) string {
	s := strings.ToLower(title)
	reg := regexp.MustCompile("[^a-z0-9]+")
	s = reg.ReplaceAllString(s, "-")
	s = strings.Trim(s, "-")
	return s
}

func parseDateInput(dateStr string) *time.Time {
	t, err := time.Parse("January 2, 2006", dateStr)
	if err == nil {
		return &t
	}
	t, err = time.Parse("2006-01-02", dateStr)
	if err == nil {
		return &t
	}
	now := time.Now()
	return &now
}

package public

import (
	"errors"
	"net/http"

	"go-cms/internal/repository"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
)

type BlogHandler struct {
	tmpl Renderer
	repo *repository.BlogRepository
}

func NewBlogHandler(tmpl Renderer, repo *repository.BlogRepository) *BlogHandler {
	return &BlogHandler{
		tmpl: tmpl,
		repo: repo,
	}
}

func (h *BlogHandler) List(w http.ResponseWriter, r *http.Request) {
	blogs, err := h.repo.FindAllPublished(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Blog — danangamw",
		"ActiveMenu": "blog",
		"Blogs":      blogs,
	}

	h.tmpl.Render(w, "blog_list", data)
}

func (h *BlogHandler) Detail(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")
	blog, err := h.repo.FindBySlug(r.Context(), slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			h.render404(w, r)
			return
		}
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if blog.Status != "published" {
		h.render404(w, r)
		return
	}

	data := map[string]any{
		"Title":      blog.Title + " — Blog",
		"ActiveMenu": "blog",
		"Blog":       blog,
	}

	h.tmpl.Render(w, "blog_detail", data)
}

func (h *BlogHandler) render404(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusNotFound)
	h.tmpl.Render(w, "404", nil)
}

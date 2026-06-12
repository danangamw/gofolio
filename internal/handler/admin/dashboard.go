package admin

import (
	"net/http"

	"go-cms/internal/middleware"
	"go-cms/internal/repository"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any)
}

type DashboardHandler struct {
	tmpl          Renderer
	blogRepo      *repository.BlogRepository
	portfolioRepo *repository.PortfolioRepository
}

func NewDashboardHandler(tmpl Renderer, blogRepo *repository.BlogRepository, portfolioRepo *repository.PortfolioRepository) *DashboardHandler {
	return &DashboardHandler{
		tmpl:          tmpl,
		blogRepo:      blogRepo,
		portfolioRepo: portfolioRepo,
	}
}

type stats struct {
	BlogCount      int64
	PortfolioCount int64
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	blogCount, err := h.blogRepo.Count(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	portfolioCount, err := h.portfolioRepo.Count(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	recent, err := h.blogRepo.Recent(r.Context(), 5)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Admin Dashboard",
		"ActiveMenu": "dashboard",
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
		"Stats": stats{
			BlogCount:      blogCount,
			PortfolioCount: portfolioCount,
		},
		"RecentBlogs": recent,
	}

	h.tmpl.Render(w, "dashboard", data)
}

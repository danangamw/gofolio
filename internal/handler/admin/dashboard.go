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
	sysConfigRepo *repository.SysConfigRepository
}

func NewDashboardHandler(tmpl Renderer, blogRepo *repository.BlogRepository, portfolioRepo *repository.PortfolioRepository, sysConfigRepo *repository.SysConfigRepository) *DashboardHandler {
	return &DashboardHandler{
		tmpl:          tmpl,
		blogRepo:      blogRepo,
		portfolioRepo: portfolioRepo,
		sysConfigRepo: sysConfigRepo,
	}
}

type stats struct {
	BlogCount      int64
	PortfolioCount int64
	SysConfigCount int64
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	blogCount, err := h.blogRepo.Count(r.Context())
	if err != nil {
		renderInternalServerError(w, h.tmpl)
		return
	}

	portfolioCount, err := h.portfolioRepo.Count(r.Context())
	if err != nil {
		renderInternalServerError(w, h.tmpl)
		return
	}

	sysConfigCount, err := h.sysConfigRepo.Count(r.Context())
	if err != nil {
		renderInternalServerError(w, h.tmpl)
		return
	}

	recent, err := h.blogRepo.Recent(r.Context(), 5)
	if err != nil {
		renderInternalServerError(w, h.tmpl)
		return
	}

	data := map[string]any{
		"Title":      "Admin Dashboard",
		"ActiveMenu": "dashboard",
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
		"Stats": stats{
			BlogCount:      blogCount,
			PortfolioCount: portfolioCount,
			SysConfigCount: sysConfigCount,
		},
		"RecentBlogs": recent,
	}

	h.tmpl.Render(w, "dashboard", data)
}

func renderInternalServerError(w http.ResponseWriter, tmpl Renderer) {
	w.WriteHeader(http.StatusInternalServerError)
	tmpl.Render(w, "500", nil)
}

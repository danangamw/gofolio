package public

import (
	"net/http"

	"go-cms/internal/repository"
)

type PortfolioHandler struct {
	tmpl Renderer
	repo *repository.PortfolioRepository
}

func NewPortfolioHandler(tmpl Renderer, repo *repository.PortfolioRepository) *PortfolioHandler {
	return &PortfolioHandler{
		tmpl: tmpl,
		repo: repo,
	}
}

func (h *PortfolioHandler) List(w http.ResponseWriter, r *http.Request) {
	portfolios, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Portfolio — danangamw",
		"ActiveMenu": "portfolio",
		"Portfolios": portfolios,
	}

	h.tmpl.Render(w, "portfolio", data)
}

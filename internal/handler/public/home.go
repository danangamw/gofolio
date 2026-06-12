package public

import (
	"net/http"

	"go-cms/internal/repository"
)

type HomeHandler struct {
	tmpl          Renderer
	blogRepo      *repository.BlogRepository
	portfolioRepo *repository.PortfolioRepository
}

func NewHomeHandler(tmpl Renderer, blogRepo *repository.BlogRepository, portfolioRepo *repository.PortfolioRepository) *HomeHandler {
	return &HomeHandler{
		tmpl:          tmpl,
		blogRepo:      blogRepo,
		portfolioRepo: portfolioRepo,
	}
}

type featuredPortfolio struct {
	Title         string
	Description   string
	Icon          string
	TechStack     []string
	ProjectURL    string
	RepositoryURL string
}

type latestBlog struct {
	Category string
	Date     string
	Title    string
	Excerpt  string
	Slug     string
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	portfolios, err := h.portfolioRepo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	blogs, err := h.blogRepo.Recent(r.Context(), 3)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Limit to top 3 portfolios on the home page
	featuredPortfolios := portfolios
	if len(featuredPortfolios) > 3 {
		featuredPortfolios = featuredPortfolios[:3]
	}

	data := map[string]any{
		"Title":              "Danang — Backend Developer Portfolio",
		"ActiveMenu":         "home",
		"FeaturedPortfolios": featuredPortfolios,
		"LatestBlogs":        blogs,
	}

	h.tmpl.Render(w, "home", data)
}

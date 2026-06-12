package admin

import (
	"net/http"
	"strings"

	"go-cms/internal/middleware"
	"go-cms/internal/model"
	"go-cms/internal/repository"

	"github.com/go-chi/chi/v5"
)

type AdminPortfolioHandler struct {
	tmpl Renderer
	repo *repository.PortfolioRepository
}

func NewAdminPortfolioHandler(tmpl Renderer, repo *repository.PortfolioRepository) *AdminPortfolioHandler {
	return &AdminPortfolioHandler{
		tmpl: tmpl,
		repo: repo,
	}
}

func (h *AdminPortfolioHandler) List(w http.ResponseWriter, r *http.Request) {
	portfolios, err := h.repo.FindAll(r.Context())
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	data := map[string]any{
		"Title":      "Manage Portfolio",
		"ActiveMenu": "portfolio",
		"Portfolios": portfolios,
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "portfolio_list_admin", data)
}

func (h *AdminPortfolioHandler) New(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "New Project",
		"ActiveMenu": "portfolio",
		"IsEdit":     false,
		"CSRFToken":  middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "portfolio_form", data)
}

func (h *AdminPortfolioHandler) Create(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	title := r.FormValue("title")
	icon := r.FormValue("icon")
	techStackStr := r.FormValue("tech_stack")
	projectURL := r.FormValue("project_url")
	repositoryURL := r.FormValue("repository_url")
	description := r.FormValue("description")

	var techStack []string
	for _, p := range strings.Split(techStackStr, ",") {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			techStack = append(techStack, trimmed)
		}
	}

	portfolio := &model.Portfolio{
		Title:         title,
		Icon:          icon,
		TechStack:     techStack,
		ProjectURL:    projectURL,
		RepositoryURL: repositoryURL,
		Description:   description,
	}

	if err := h.repo.Create(r.Context(), portfolio); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/portfolios", http.StatusSeeOther)
}

func (h *AdminPortfolioHandler) Edit(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	portfolio, err := h.repo.FindByTitle(r.Context(), title)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if portfolio == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	data := map[string]any{
		"Title":           "Edit Project",
		"ActiveMenu":      "portfolio",
		"IsEdit":          true,
		"Portfolio":       portfolio,
		"TechStackString": strings.Join(portfolio.TechStack, ", "),
		"CSRFToken":       middleware.GetCSRFToken(r.Context()),
	}
	h.tmpl.Render(w, "portfolio_form", data)
}

func (h *AdminPortfolioHandler) Update(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	portfolio, err := h.repo.FindByTitle(r.Context(), title)
	if err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	if portfolio == nil {
		http.Error(w, "Not Found", http.StatusNotFound)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Bad Request", http.StatusBadRequest)
		return
	}

	newTitle := r.FormValue("title")
	icon := r.FormValue("icon")
	techStackStr := r.FormValue("tech_stack")
	projectURL := r.FormValue("project_url")
	repositoryURL := r.FormValue("repository_url")
	description := r.FormValue("description")

	var techStack []string
	for _, p := range strings.Split(techStackStr, ",") {
		trimmed := strings.TrimSpace(p)
		if trimmed != "" {
			techStack = append(techStack, trimmed)
		}
	}

	portfolio.Title = newTitle
	portfolio.Icon = icon
	portfolio.TechStack = techStack
	portfolio.ProjectURL = projectURL
	portfolio.RepositoryURL = repositoryURL
	portfolio.Description = description

	if err := h.repo.Update(r.Context(), portfolio); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/admin/portfolios", http.StatusSeeOther)
}

func (h *AdminPortfolioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	if err := h.repo.Delete(r.Context(), title); err != nil {
		http.Error(w, "Internal Server Error: "+err.Error(), http.StatusInternalServerError)
		return
	}
	http.Redirect(w, r, "/admin/portfolios", http.StatusSeeOther)
}

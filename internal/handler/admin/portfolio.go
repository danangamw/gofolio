package admin

import (
	"net/http"
	"strings"

	"github.com/go-chi/chi/v5"
)

type AdminPortfolioHandler struct {
	tmpl Renderer
}

func NewAdminPortfolioHandler(tmpl Renderer) *AdminPortfolioHandler {
	return &AdminPortfolioHandler{tmpl: tmpl}
}

type adminPortfolioItem struct {
	Title         string
	Description   string
	Icon          string
	TechStack     []string
	ProjectURL    string
	RepositoryURL string
}

var mockPortfolios = []adminPortfolioItem{
	{
		Title:         "Go-CMS Monolith",
		Description:   "A modular, high-performance Content Management System built with Go, Postgres, Redis, and OpenTelemetry.",
		Icon:          "🚀",
		TechStack:     []string{"Go", "PostgreSQL", "Redis", "OTel"},
		ProjectURL:    "https://github.com",
		RepositoryURL: "https://github.com",
	},
	{
		Title:         "E-Commerce Payment Gateway",
		Description:   "An asynchronous microservices system to process e-commerce transactions using Message Queues.",
		Icon:          "💳",
		TechStack:     []string{"Go", "RabbitMQ", "PostgreSQL", "Docker"},
		ProjectURL:    "",
		RepositoryURL: "https://github.com",
	},
	{
		Title:         "Realtime Chat Application",
		Description:   "A real-time communication application using WebSockets, distributed with Redis Pub/Sub.",
		Icon:          "💬",
		TechStack:     []string{"Go", "WebSocket", "Redis", "HTML/CSS"},
		ProjectURL:    "https://github.com",
		RepositoryURL: "https://github.com",
	},
}

func (h *AdminPortfolioHandler) List(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "Manage Portfolio",
		"ActiveMenu": "portfolio",
		"Portfolios": mockPortfolios,
	}
	h.tmpl.Render(w, "portfolio_list_admin", data)
}

func (h *AdminPortfolioHandler) New(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "New Project",
		"ActiveMenu": "portfolio",
		"IsEdit":     false,
	}
	h.tmpl.Render(w, "portfolio_form", data)
}

func (h *AdminPortfolioHandler) Create(w http.ResponseWriter, r *http.Request) {
	// Simple redirect back to list for static demo
	http.Redirect(w, r, "/admin/portfolios", http.StatusSeeOther)
}

func (h *AdminPortfolioHandler) Edit(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	var found adminPortfolioItem
	for _, p := range mockPortfolios {
		if p.Title == title {
			found = p
			break
		}
	}
	if found.Title == "" {
		found = mockPortfolios[0] // Fallback
	}

	data := map[string]any{
		"Title":           "Edit Project",
		"ActiveMenu":      "portfolio",
		"IsEdit":          true,
		"Portfolio":       found,
		"TechStackString": strings.Join(found.TechStack, ", "),
	}
	h.tmpl.Render(w, "portfolio_form", data)
}

func (h *AdminPortfolioHandler) Update(w http.ResponseWriter, r *http.Request) {
	// Simple redirect back to list for static demo
	http.Redirect(w, r, "/admin/portfolios", http.StatusSeeOther)
}

func (h *AdminPortfolioHandler) Delete(w http.ResponseWriter, r *http.Request) {
	title := chi.URLParam(r, "title")
	var newPortfolios []adminPortfolioItem
	for _, p := range mockPortfolios {
		if p.Title != title {
			newPortfolios = append(newPortfolios, p)
		}
	}
	if len(newPortfolios) < len(mockPortfolios) {
		mockPortfolios = newPortfolios
	}
	http.Redirect(w, r, "/admin/portfolios", http.StatusSeeOther)
}

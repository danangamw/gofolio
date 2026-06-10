package public

import (
	"net/http"
)

type HomeHandler struct {
	tmpl Renderer
}

func NewHomeHandler(tmpl Renderer) *HomeHandler {
	return &HomeHandler{tmpl: tmpl}
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
	featuredPortfolios := []featuredPortfolio{
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

	latestBlogs := []latestBlog{
		{
			Category: "Go",
			Date:     "June 10, 2026",
			Title:    "Understanding Dependency Injection in Go Simply",
			Excerpt:  "How to neatly manage database and third-party dependencies in Go applications without external frameworks.",
			Slug:     "understanding-dependency-injection-in-go",
		},
		{
			Category: "Observability",
			Date:     "June 8, 2026",
			Title:    "Integrating OpenTelemetry Tracing in GORM",
			Excerpt:  "A complete guide on recording SQL query performance directly to Grafana Tempo using the OTel GORM plugin.",
			Slug:     "integrating-opentelemetry-tracing-in-gorm",
		},
		{
			Category: "Database",
			Date:     "June 5, 2026",
			Title:    "Designing Database Schema Migrations with Atlas",
			Excerpt:  "Why declarative migrations with Atlas are safer and more efficient than traditional manual SQL scripts.",
			Slug:     "designing-database-schema-migrations-with-atlas",
		},
	}

	data := map[string]any{
		"Title":              "Danang — Backend Developer Portfolio",
		"ActiveMenu":         "home",
		"FeaturedPortfolios": featuredPortfolios,
		"LatestBlogs":        latestBlogs,
	}

	h.tmpl.Render(w, "home", data)
}

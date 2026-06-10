package public

import (
	"net/http"
)

type PortfolioHandler struct {
	tmpl Renderer
}

func NewPortfolioHandler(tmpl Renderer) *PortfolioHandler {
	return &PortfolioHandler{tmpl: tmpl}
}

type portfolioItem struct {
	Title         string
	Description   string
	Icon          string
	TechStack     []string
	ProjectURL    string
	RepositoryURL string
}

func (h *PortfolioHandler) List(w http.ResponseWriter, r *http.Request) {
	portfolios := []portfolioItem{
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
		{
			Title:         "Docker Janitor CLI",
			Description:   "A lightweight CLI application to automatically clean up Docker containers, images, volumes, and networks.",
			Icon:          "🧹",
			TechStack:     []string{"Go", "Docker API", "CLI"},
			ProjectURL:    "",
			RepositoryURL: "https://github.com/danangamw/go-janitor",
		},
		{
			Title:         "Notification Engine",
			Description:   "A mass notification engine for email, WhatsApp, and push notifications with automatic failover.",
			Icon:          "🔔",
			TechStack:     []string{"Go", "gRPC", "PostgreSQL", "Redis"},
			ProjectURL:    "",
			RepositoryURL: "https://github.com",
		},
	}

	data := map[string]any{
		"Title":      "Portfolio — Danang",
		"ActiveMenu": "portfolio",
		"Portfolios": portfolios,
	}

	h.tmpl.Render(w, "portfolio", data)
}

package admin

import "net/http"

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any)
}

type DashboardHandler struct {
	tmpl Renderer
}

func NewDashboardHandler(tmpl Renderer) *DashboardHandler {
	return &DashboardHandler{tmpl: tmpl}
}

type recentBlog struct {
	Title    string
	Category string
	Date     string
	Slug     string
}

type stats struct {
	BlogCount      int
	PortfolioCount int
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	recent := []recentBlog{
		{
			Title:    "Understanding Dependency Injection in Go Simply",
			Category: "Go",
			Date:     "June 10, 2026",
			Slug:     "understanding-dependency-injection-in-go",
		},
		{
			Title:    "Integrating OpenTelemetry Tracing in GORM",
			Category: "Observability",
			Date:     "June 8, 2026",
			Slug:     "integrating-opentelemetry-tracing-in-gorm",
		},
	}

	data := map[string]any{
		"Title":      "Admin Dashboard",
		"ActiveMenu": "dashboard",
		"Stats": stats{
			BlogCount:      12,
			PortfolioCount: 5,
		},
		"RecentBlogs": recent,
	}

	h.tmpl.Render(w, "dashboard", data)
}

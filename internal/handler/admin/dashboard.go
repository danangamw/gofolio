package admin

import "net/http"

type DashboardHandler struct{}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{}
}

func (h *DashboardHandler) Index(w http.ResponseWriter, r *http.Request) {
	// TODO: render dashboard.html with blog/portfolio counts
}

package public

import "net/http"

type PortfolioHandler struct{}

func NewPortfolioHandler() *PortfolioHandler {
	return &PortfolioHandler{}
}

func (h *PortfolioHandler) List(w http.ResponseWriter, r *http.Request) {
	// TODO: render portfolio.html template
}

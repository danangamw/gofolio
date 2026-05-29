package public

import "net/http"

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	// TODO: render home.html template
}

package public

import "net/http"

type AboutHandler struct{}

func NewAboutHandler() *AboutHandler {
	return &AboutHandler{}
}

func (h *AboutHandler) Index(w http.ResponseWriter, r *http.Request) {
	// TODO: render about.html template
}

package public

import (
	"net/http"
)

type Renderer interface {
	Render(w http.ResponseWriter, name string, data any)
}

type AboutHandler struct {
	tmpl Renderer
}

func NewAboutHandler(tmpl Renderer) *AboutHandler {
	return &AboutHandler{tmpl: tmpl}
}

func (h *AboutHandler) Index(w http.ResponseWriter, r *http.Request) {
	data := map[string]any{
		"Title":      "Tentang Saya — Danang",
		"ActiveMenu": "about",
	}

	h.tmpl.Render(w, "about", data)
}

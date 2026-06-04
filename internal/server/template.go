package server

import (
	"html/template"
	"net/http"
)

// loadTemplate parse all template from embed.FS
// call once on startup

func (s *Server) loadTemplates() (*template.Template, error) {
	// parse layout first as base
	tmpl, err := template.ParseFS(s.webFs,
		"web/templates/layout/base.html",
		"web/templates/layout/admin_base.html",
	)

	if err != nil {
		return nil, err
	}

	// parse all page template
	pages := []string{
		"web/templates/public/*.html",
		"web/templates/auth/*.html",
		"web/templates/admin/*.html",
	}

	for _, page := range pages {
		tmpl, err = tmpl.ParseFS(s.webFs, page)
		if err != nil {
			return nil, err
		}
	}

	return tmpl, nil
}

// render helper to render template + handle error
func (s *Server) render(w http.ResponseWriter, tmpl *template.Template, name string, data any) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if err := tmpl.ExecuteTemplate(w, name, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

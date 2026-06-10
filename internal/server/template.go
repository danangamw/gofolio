package server

import (
	"html/template"
	"net/http"
)

// TemplateRegistry holds a map of compiled page templates to prevent block name collisions.
type TemplateRegistry struct {
	templates map[string]*template.Template
}

// Render renders a specific page using its registered template and matches the appropriate layout.
func (r *TemplateRegistry) Render(w http.ResponseWriter, name string, data any) {
	tmpl, ok := r.templates[name]
	if !ok {
		http.Error(w, "Template not found: "+name, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	// Determine layout dynamically based on lookups.
	layout := "base"
	if tmpl.Lookup("admin_base") != nil {
		layout = "admin_base"
	}

	if err := tmpl.ExecuteTemplate(w, layout, data); err != nil {
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
}

func (s *Server) loadTemplates() (*TemplateRegistry, error) {
	registry := &TemplateRegistry{
		templates: make(map[string]*template.Template),
	}

	// Define public and auth pages that use layout/base.html
	basePages := []struct {
		name string
		path string
	}{
		{"home", "web/templates/public/home.html"},
		{"portfolio", "web/templates/public/portfolio.html"},
		{"blog_list", "web/templates/public/blog_list.html"},
		{"blog_detail", "web/templates/public/blog_detail.html"},
		{"about", "web/templates/public/about.html"},
		{"login", "web/templates/auth/login.html"},
	}

	for _, p := range basePages {
		tmpl, err := template.ParseFS(s.webFs, "web/templates/layout/base.html", p.path)
		if err != nil {
			return nil, err
		}
		registry.templates[p.name] = tmpl
	}

	// Define admin pages that use layout/admin_base.html
	adminPages := []struct {
		name string
		path string
	}{
		{"dashboard", "web/templates/admin/dashboard.html"},
		{"blog_list_admin", "web/templates/admin/blog_list.html"},
		{"blog_form", "web/templates/admin/blog_form.html"},
		{"portfolio_list_admin", "web/templates/admin/portfolio_list.html"},
		{"portfolio_form", "web/templates/admin/portfolio_form.html"},
	}

	for _, p := range adminPages {
		tmpl, err := template.ParseFS(s.webFs, "web/templates/layout/admin_base.html", p.path)
		if err != nil {
			return nil, err
		}
		registry.templates[p.name] = tmpl
	}

	// Define error pages that use layout/base.html
	errorPages := []string{"404", "500"}
	for _, p := range errorPages {
		tmpl, err := template.ParseFS(s.webFs, "web/templates/layout/base.html", "web/templates/public/"+p+".html")
		if err != nil {
			return nil, err
		}
		registry.templates[p] = tmpl
	}

	return registry, nil
}

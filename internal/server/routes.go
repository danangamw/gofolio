package server

import (
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	adminhandler "go-cms/internal/handler/admin"
	authhandler "go-cms/internal/handler/auth"
	publichandler "go-cms/internal/handler/public"
	cmsmiddleware "go-cms/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	// Load all templates once at startup — panic if templates are broken.
	tmpl, err := s.loadTemplates()
	if err != nil {
		panic("failed to load templates: " + err.Error())
	}
	s.tmpl = tmpl

	r := chi.NewRouter()

	// ── Global middleware ────────────────────────────────────────────────────
	r.Use(middleware.Recoverer)
	r.Use(middleware.RequestID)
	r.Use(cmsmiddleware.OtelHTTP(s.cfg.ServiceName))
	r.Use(cmsmiddleware.TraceIDHeader)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Trace-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	// ── Custom Not Found Handler ─────────────────────────────────────────────
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusNotFound)
		tmpl.Render(w, "404", nil)
	})

	// ── Static assets from embed.FS ──────────────────────────────────────────
	webSub, err := fs.Sub(s.webFs, "web")
	if err != nil {
		panic("failed to get web sub filesystem: " + err.Error())
	}
	r.Handle("/static/*", http.FileServer(http.FS(webSub)))

	// ── Health ───────────────────────────────────────────────────────────────
	r.Get("/health", s.healthHandler)

	// ── Public Pages ─────────────────────────────────────────────────────────
	homeH := publichandler.NewHomeHandler(tmpl)
	portH := publichandler.NewPortfolioHandler(tmpl)
	blogH := publichandler.NewBlogHandler(tmpl)
	aboutH := publichandler.NewAboutHandler(tmpl)

	r.Get("/", homeH.Index)
	r.Get("/portfolio", portH.List)
	r.Get("/blog", blogH.List)
	r.Get("/blog/{slug}", blogH.Detail)
	r.Get("/about", aboutH.Index)

	// ── Auth handler ─────────────────────────────────────────────────────────
	authH := authhandler.New(tmpl, s.userRepo, s.sessions)

	r.Get("/login", authH.LoginPage)
	r.Post("/login", authH.Login)

	// ── Admin Pages ──────────────────────────────────────────────────────────
	adminDashH := adminhandler.NewDashboardHandler(tmpl)
	adminBlogH := adminhandler.NewAdminBlogHandler(tmpl)
	adminPortH := adminhandler.NewAdminPortfolioHandler(tmpl)

	// ── Admin routes (bypassed login for static testing) ──────────────────────
	r.Route("/admin", func(r chi.Router) {
		// r.Use(cmsmiddleware.Auth(s.sessions)) // Bypassed for static UI testing

		r.Post("/logout", authH.Logout)

		r.Get("/", adminDashH.Index)

		// Blog administration routes
		r.Get("/blogs", adminBlogH.List)
		r.Get("/blogs/new", adminBlogH.New)
		r.Post("/blogs/new", adminBlogH.Create)
		r.Get("/blogs/edit/{slug}", adminBlogH.Edit)
		r.Post("/blogs/edit/{slug}", adminBlogH.Update)
		r.Post("/blogs/delete/{slug}", adminBlogH.Delete)

		// Portfolio administration routes
		r.Get("/portfolios", adminPortH.List)
		r.Get("/portfolios/new", adminPortH.New)
		r.Post("/portfolios/new", adminPortH.Create)
		r.Get("/portfolios/edit/{title}", adminPortH.Edit)
		r.Post("/portfolios/edit/{title}", adminPortH.Update)
		r.Post("/portfolios/delete/{title}", adminPortH.Delete)
	})

	return r
}

// ── Handler implementations ───────────────────────────────────────────────────

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

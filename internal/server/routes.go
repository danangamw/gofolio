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
	"go-cms/internal/service"
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
	r.Use(cmsmiddleware.CSRF)

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

	// Serve uploads locally if configured
	if s.cfg.UploadStorage == "local" {
		r.Handle("/uploads/*", http.StripPrefix("/uploads/", http.FileServer(http.Dir(s.cfg.UploadDir))))
	}

	// ── Health ───────────────────────────────────────────────────────────────
	r.Get("/health", s.healthHandler)

	// ── Public Pages ─────────────────────────────────────────────────────────
	homeH := publichandler.NewHomeHandler(tmpl, s.blogRepo, s.portfolioRepo)
	portH := publichandler.NewPortfolioHandler(tmpl, s.portfolioRepo)
	blogH := publichandler.NewBlogHandler(tmpl, s.blogRepo)
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

	// ── Admin Pages & Services ───────────────────────────────────────────────
	uploadSvc := service.NewUploadService(s.cfg)
	uploadH := adminhandler.NewUploadHandler(uploadSvc)

	adminDashH := adminhandler.NewDashboardHandler(tmpl, s.blogRepo, s.portfolioRepo)
	adminBlogH := adminhandler.NewAdminBlogHandler(tmpl, s.blogRepo)
	adminPortH := adminhandler.NewAdminPortfolioHandler(tmpl, s.portfolioRepo)

	// ── Admin routes ─────────────────────────────────────────────────────────
	r.Route("/admin", func(r chi.Router) {
		r.Use(cmsmiddleware.Auth(s.sessions))

		r.Post("/logout", authH.Logout)

		r.Get("/", adminDashH.Index)

		// File Upload route
		r.Post("/upload", uploadH.Upload)

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

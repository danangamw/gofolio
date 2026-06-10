package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	authhandler "go-cms/internal/handler/auth"
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

	// ── Static assets from embed.FS ──────────────────────────────────────────
	r.Handle("/static/*", http.FileServer(http.FS(s.webFs)))

	// ── Health ───────────────────────────────────────────────────────────────
	r.Get("/health", s.healthHandler)

	// ── Auth handler ─────────────────────────────────────────────────────────
	authH := authhandler.New(tmpl, s.userRepo, s.sessions)

	r.Get("/login", authH.LoginPage)
	r.Post("/login", authH.Login)

	// ── Admin routes (protected) ──────────────────────────────────────────────
	r.Route("/admin", func(r chi.Router) {
		r.Use(cmsmiddleware.Auth(s.sessions))

		r.Post("/logout", authH.Logout)

		// Dashboard placeholder — will be replaced in Phase 4.
		r.Get("/", s.adminDashboardHandler)
	})

	return r
}

// ── Handler implementations ───────────────────────────────────────────────────

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

// adminDashboardHandler is a placeholder until Phase 4 (Admin Dashboard).
func (s *Server) adminDashboardHandler(w http.ResponseWriter, r *http.Request) {
	userID := cmsmiddleware.GetUserID(r)
	w.Header().Set("Content-Type", "text/plain")
	log.Printf("admin dashboard accessed by userID=%s", userID)
	w.Write([]byte("Welcome to the admin dashboard. (Phase 4 will render the full UI)"))
}

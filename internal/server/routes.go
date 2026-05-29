package server

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	cmsmiddleware "go-cms/internal/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := chi.NewRouter()

	// Recoverer catches panics and returns 500 without crashing
	r.Use(middleware.Recoverer)
	// RequestID injects a unique ID into every request for internal tracing
	r.Use(middleware.RequestID)
	// OtelHTTP: auto-instrument spans + metrics for each request
	r.Use(cmsmiddleware.OtelHTTP(s.serviceName))
	// Inject X-Trace-ID into the response header so the client can trace it in Grafana
	r.Use(cmsmiddleware.TraceIDHeader)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS", "PATCH"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Trace-ID"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	r.Get("/", s.HelloWorldHandler)

	r.Get("/health", s.healthHandler)

	return r
}

func (s *Server) HelloWorldHandler(w http.ResponseWriter, r *http.Request) {
	resp := make(map[string]string)
	resp["message"] = "Hello World"

	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("error handling JSON marshal. Err: %v", err)
	}

	_, _ = w.Write(jsonResp)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	jsonResp, _ := json.Marshal(s.db.Health())
	_, _ = w.Write(jsonResp)
}

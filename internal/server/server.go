package server

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-cms/internal/config"
	"go-cms/internal/database"
)

type Server struct {
	cfg   *config.Config
	db    database.Service
	webFs embed.FS
}

func NewServer(cfg *config.Config, db database.Service, webFs embed.FS) *http.Server {
	s := &Server{
		cfg:   cfg,
		db:    db,
		webFs: webFs,
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

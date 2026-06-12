package server

import (
	"embed"
	"fmt"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-cms/internal/config"
	"go-cms/internal/database"
	"go-cms/internal/repository"
	"go-cms/internal/session"
)

// Server holds all application-wide dependencies shared across handlers.
type Server struct {
	cfg           *config.Config
	db            database.Service
	webFs         embed.FS
	sessions      session.Store
	userRepo      *repository.UserRepository
	blogRepo      *repository.BlogRepository
	portfolioRepo *repository.PortfolioRepository
	tmpl          *TemplateRegistry // parsed once at startup
}

// NewServer wires all dependencies and returns a ready *http.Server.
func NewServer(
	cfg *config.Config,
	db database.Service,
	webFs embed.FS,
	sessions session.Store,
) *http.Server {
	s := &Server{
		cfg:           cfg,
		db:            db,
		webFs:         webFs,
		sessions:      sessions,
		userRepo:      repository.NewUserRepository(db.GetDB()),
		blogRepo:      repository.NewBlogRepository(db.GetDB()),
		portfolioRepo: repository.NewPortfolioRepository(db.GetDB()),
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

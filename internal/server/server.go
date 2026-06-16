package server

import (
	"context"
	"embed"
	"fmt"
	"log"
	"net/http"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-cms/internal/config"
	"go-cms/internal/database"
	"go-cms/internal/repository"
	"go-cms/internal/session"
	"go-cms/pkg/storage"
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
	sysConfigRepo *repository.SysConfigRepository
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
		sysConfigRepo: repository.NewSysConfigRepository(db.GetDB()),
	}

	if err := s.sysConfigRepo.SeedDefaults(context.Background()); err != nil {
		log.Printf("WARN: failed to seed system configs (continuing with fallback content): %v", err)
	}

	if cfg.UploadStorage == "s3" {
		log.Printf("Initializing S3/MinIO storage client (endpoint: %s, bucket: %s)...", cfg.S3Endpoint, cfg.S3Bucket)
		err := storage.InitStorageFromConfig(
			cfg.S3Endpoint,
			cfg.S3PublicEndpoint,
			cfg.S3AccessKeyID,
			cfg.S3SecretAccessKey,
			cfg.S3Bucket,
			cfg.S3Region,
		)
		if err != nil {
			log.Fatalf("Failed to initialize S3/MinIO storage: %v", err)
		}
		log.Println("S3/MinIO storage client initialized successfully.")
	}

	return &http.Server{
		Addr:         fmt.Sprintf(":%s", cfg.AppPort),
		Handler:      s.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
}

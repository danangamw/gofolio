package main

import (
	"context"
	"log"
	"time"

	"go-cms/internal/config"
	"go-cms/internal/database"
	"go-cms/internal/model"
	"go-cms/internal/repository"
	"go-cms/internal/service"
)

func main() {
	cfg := config.Load()

	db := database.New(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepository(db.GetDB())

	// Force GORM auto-migrations to ensure tables exist
	log.Println("Running GORM auto-migrations for seeder...")
	if err := db.GetDB().AutoMigrate(&model.User{}, &model.Blog{}, &model.Portfolio{}, &model.Session{}); err != nil {
		log.Fatalf("seed: auto migration failed: %v", err)
	}
	log.Println("GORM auto-migrations completed successfully.")

	username := cfg.AdminUsername
	password := cfg.AdminPassword

	if username == "" || password == "" {
		log.Fatal("ADMIN_USERNAME and ADMIN_PASSWORD must be set in .env")
	}

	ctx := context.Background()

	// Check if user already exists — idempotent.
	existing, err := userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Fatalf("seed: check existing user: %v", err)
	}
	if existing != nil {
		log.Printf("Admin user %q already exists — skipping seed.", username)
		return
	}

	hash, err := service.HashPassword(password)
	if err != nil {
		log.Fatalf("seed: hash password: %v", err)
	}

	user := &model.User{
		Username:     username,
		PasswordHash: hash,
	}

	if err := userRepo.Create(ctx, user); err != nil {
		log.Fatalf("seed: create admin user: %v", err)
	}

	log.Printf("Admin user %q created successfully (id: %s)", username, user.ID)

	// Seed Blog Posts
	blogRepo := repository.NewBlogRepository(db.GetDB())
	count, err := blogRepo.Count(ctx)
	if err == nil && count == 0 {
		now := time.Now()
		blogs := []model.Blog{
			{
				Title:       "Understanding Dependency Injection in Go Simply",
				Slug:        "understanding-dependency-injection-in-go",
				Category:    "Go",
				Excerpt:     "How to neatly manage database and third-party dependencies in Go applications without external frameworks.",
				Content:     `Dependency Injection (DI) is often considered a complex concept because it is associated with large frameworks. However, in the Go programming language, Dependency Injection is actually very simple and does not require additional frameworks (like Wire or Dig) for most applications.

## What is Dependency Injection?
Simply put, Dependency Injection means we pass the dependencies (like database connections or third-party clients) required by a function/struct, instead of letting it instantiate or search for them itself from global variables.

## Without DI Approach (Bad)
` + "```go" + `
package handler

import "database/sql"

var DB *sql.DB // Global variable

func GetUserHandler(w http.ResponseWriter, r *http.Request) {
    // Tightly coupled to a global database variable
    rows, err := DB.Query("SELECT id, username FROM users")
    // ...
}
` + "```" + `
This approach makes unit testing difficult because the handler is tightly coupled to a global database variable.

## With DI Approach (Good)
The best way in Go is to create a handler struct that accepts a database interface during initialization:
` + "```go" + `
type UserHandler struct {
    db *sql.DB
}

func NewUserHandler(db *sql.DB) *UserHandler {
    return &UserHandler{db: db}
}

func (h *UserHandler) GetUsers(w http.ResponseWriter, r *http.Request) {
    // Use h.db here
}
` + "```" + `

## Conclusion
By implementing simple Dependency Injection through struct constructors (like NewUserHandler), your Go application code becomes much easier to test, flexible, and exceptionally modular.`,
				Status:      "published",
				PublishedAt: &now,
			},
			{
				Title:       "Integrating OpenTelemetry Tracing in GORM",
				Slug:        "integrating-opentelemetry-tracing-in-gorm",
				Category:    "Observability",
				Excerpt:     "A complete guide on recording SQL query performance directly to Grafana Tempo using the OTel GORM plugin.",
				Content:     `When our application slows down, one of the main suspects is sub-optimal database queries. By setting up distributed tracing using OpenTelemetry (OTel) and GORM, we can track SQL query details, parameters, and duration directly in a visual dashboard like Grafana Tempo.

## Why Use Tracing for Databases?
With tracing, every database query executed during an HTTP request will be recorded as a child span. This makes it easy to see the relation between HTTP requests and the SQL queries running under the hood.

## Integration Steps
First, install the GORM tracing plugin:
` + "```bash" + `
go get gorm.io/plugin/opentelemetry/tracing
` + "```" + `
Then, after opening the GORM connection, register the plugin:
` + "```go" + `
import (
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
    "gorm.io/plugin/opentelemetry/tracing"
)

db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
if err != nil {
    log.Fatal(err)
}

// Connect GORM with OpenTelemetry
if err := db.Use(tracing.NewPlugin(tracing.WithoutMetrics())); err != nil {
    log.Printf("Failed to register OTel plugin: %v", err)
}
` + "```" + `

## Importance of Context
For tracing to work, you must always pass the request context (ctx) when executing database queries in GORM. Example:
` + "```go" + `
// Span context is passed down via WithContext
err := db.WithContext(r.Context()).Where("id = ?", id).First(&user).Error
` + "```" + `

## Conclusion
Integrating OpenTelemetry with GORM is a crucial step for production-ready applications. It saves hours of debugging time when database performance issues arise.`,
				Status:      "published",
				PublishedAt: &now,
			},
			{
				Title:       "Designing Database Schema Migrations with Atlas",
				Slug:        "designing-database-schema-migrations-with-atlas",
				Category:    "Database",
				Excerpt:     "Why declarative migrations with Atlas are safer and more efficient than traditional manual SQL scripts.",
				Content:     `Managing database schema changes (migrations) often poses a major challenge when working in teams. Using manual approaches like writing raw SQL files is prone to conflicts and human errors. This is where Atlas shines as a modern migration solution.

## Declarative vs. Imperative Approach
Most traditional Go migration libraries (like golang-migrate) use an imperative model: you write explicit SQL commands like CREATE TABLE or ALTER TABLE. If there is a typo, the migration could fail halfway.

Atlas uses a declarative approach: you describe the desired final state of the database (e.g., from GORM structs), and Atlas automatically calculates the safe transition SQL (diff) required.

## Integrating GORM with Atlas
Atlas can read GORM structs directly using a custom loader, comparing them to a local database to generate migration files automatically.
` + "```bash" + `
# Command to generate a new migration file
atlas migrate diff migration_name --env local
` + "```" + `

## Conclusion
By combining declarative safety from Atlas and GORM mapping convenience, developers can change database schemas confidently, with minimal risk, and integrate them smoothly into CI/CD pipelines.`,
				Status:      "published",
				PublishedAt: &now,
			},
		}

		for _, b := range blogs {
			if err := blogRepo.Create(ctx, &b); err != nil {
				log.Printf("seed: error creating blog %q: %v", b.Title, err)
			}
		}
		log.Println("Seeded initial blog posts successfully.")
	}

	// Seed Portfolios
	portfolioRepo := repository.NewPortfolioRepository(db.GetDB())
	pCount, err := portfolioRepo.Count(ctx)
	if err == nil && pCount == 0 {
		portfolios := []model.Portfolio{
			{
				Title:         "Go-CMS Monolith",
				Icon:          "🚀",
				Description:   "A modular, high-performance Content Management System built with Go, Postgres, Redis, and OpenTelemetry.",
				TechStack:     []string{"Go", "PostgreSQL", "Redis", "OTel"},
				ProjectURL:    "https://github.com",
				RepositoryURL: "https://github.com",
			},
			{
				Title:         "E-Commerce Payment Gateway",
				Icon:          "💳",
				Description:   "An asynchronous microservices system to process e-commerce transactions using Message Queues.",
				TechStack:     []string{"Go", "RabbitMQ", "PostgreSQL", "Docker"},
				ProjectURL:    "",
				RepositoryURL: "https://github.com",
			},
			{
				Title:         "Realtime Chat Application",
				Icon:          "💬",
				Description:   "A real-time communication application using WebSockets, distributed with Redis Pub/Sub.",
				TechStack:     []string{"Go", "WebSocket", "Redis", "HTML/CSS"},
				ProjectURL:    "https://github.com",
				RepositoryURL: "https://github.com",
			},
			{
				Title:         "Docker Janitor CLI",
				Icon:          "🧹",
				Description:   "A lightweight CLI application to automatically clean up Docker containers, images, volumes, and networks.",
				TechStack:     []string{"Go", "Docker API", "CLI"},
				ProjectURL:    "",
				RepositoryURL: "https://github.com/danangamw/go-janitor",
			},
			{
				Title:         "Notification Engine",
				Icon:          "🔔",
				Description:   "A mass notification engine for email, WhatsApp, and push notifications with automatic failover.",
				TechStack:     []string{"Go", "gRPC", "PostgreSQL", "Redis"},
				ProjectURL:    "",
				RepositoryURL: "https://github.com",
			},
		}

		for _, p := range portfolios {
			if err := portfolioRepo.Create(ctx, &p); err != nil {
				log.Printf("seed: error creating portfolio %q: %v", p.Title, err)
			}
		}
		log.Println("Seeded initial portfolios successfully.")
	}
}

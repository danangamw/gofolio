package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os/signal"
	"syscall"
	"time"

	gocms "go-cms"
	"go-cms/internal/config"
	"go-cms/internal/database"
	"go-cms/internal/server"
	"go-cms/internal/session"
	"go-cms/internal/telemetry"
	"go-cms/pkg/logger"
)

func gracefulShutdown(apiServer *http.Server, done chan bool) {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	<-ctx.Done()

	log.Println("shutting down gracefully, press Ctrl+C again to force")
	stop()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := apiServer.Shutdown(ctx); err != nil {
		log.Printf("Server forced to shutdown with error: %v", err)
	}

	log.Println("Server exiting")
	done <- true
}

// webFS is defined in the root package (web.go) next to the web/ directory.

func main() {
	ctx := context.Background()
	cfg := config.Load()

	// Initialize Telemetry
	shutdown, err := telemetry.Init(ctx, telemetry.Config{
		ServiceName:    cfg.ServiceName,
		ServiceVersion: cfg.ServiceVersion,
		Environment:    cfg.AppEnv,
		CollectorAddr:  cfg.OTLPEndpoint,
	})
	if err != nil {
		log.Printf("WARN: telemetry initialization failed (proceeding without it): %v", err)
		shutdown = func(context.Context) error { return nil }
	}
	defer func() {
		shutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		if err := shutdown(shutCtx); err != nil {
			log.Printf("telemetry shutdown error: %v", err)
		}
	}()

	// Initialize Logger
	logger.Init(logger.Options{
		ServiceName: cfg.ServiceName,
		Level:       cfg.LogLevel,
		Environment: cfg.AppEnv,
	})

	db := database.New(cfg)
	defer db.Close()

	// Init session store — Redis if REDIS_URL is set, Postgres otherwise.
	sessions, err := session.NewStore(cfg.RedisURL, db)
	if err != nil {
		log.Fatalf("Failed to init session store: %v", err)
	}

	srv := server.NewServer(cfg, db, gocms.WebFS, sessions)

	done := make(chan bool, 1)

	go gracefulShutdown(srv, done)

	logger.Info(ctx, fmt.Sprintf("Server running on :%s (env: %s)", cfg.AppPort, cfg.AppEnv))
	err = srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

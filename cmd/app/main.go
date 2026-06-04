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

	cfg := config.Load()

	db := database.New(cfg)
	defer db.Close()

	srv := server.NewServer(cfg, db, gocms.WebFS)

	done := make(chan bool, 1)

	go gracefulShutdown(srv, done)

	logger.Info(context.Background(), "Server running on :%s (env: %s)", cfg.AppPort, cfg.AppEnv)
	err := srv.ListenAndServe()
	if err != nil && err != http.ErrServerClosed {
		panic(fmt.Sprintf("server error: %s", err))
	}

	<-done
	log.Println("Graceful shutdown complete.")
}

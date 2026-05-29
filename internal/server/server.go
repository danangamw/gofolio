package server

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"

	_ "github.com/joho/godotenv/autoload"

	"go-cms/internal/config"
	"go-cms/internal/database"
)

type Server struct {
	port        int
	serviceName string

	db database.Service
}

func NewServer() *http.Server {
	cfg := config.Load()
	port, _ := strconv.Atoi(os.Getenv("APP_PORT"))
	NewServer := &Server{
		port:        port,
		serviceName: cfg.ServiceName,

		db: database.New(),
	}

	// Declare Server config
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", NewServer.port),
		Handler:      NewServer.RegisterRoutes(),
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	return server
}

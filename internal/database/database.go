package database

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"
	"time"

	"go-cms/internal/model"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// Service represents a service that interacts with a database.
type Service interface {
	// Health returns a map of health status information.
	Health() map[string]string

	// Close terminates the database connection.
	Close() error

	// GetDB returns the active GORM DB instance.
	GetDB() *gorm.DB
}

type service struct {
	db *gorm.DB
}

var (
	dbInstance  *service
)

func New() Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		log.Fatal("DATABASE_URL environment variable is required but not set")
	}

	// Determine GORM log level based on APP_ENV
	gormLogLevel := logger.Silent
	if os.Getenv("APP_ENV") == "development" {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(connStr), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve generic database object: %v", err)
	}

	// Connection Pooling Settings
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully with GORM")

	// Run GORM AutoMigrate if configured (Optional, usually for dev)
	if os.Getenv("APP_AUTO_MIGRATE") == "true" {
		log.Println("Running GORM auto-migrations...")
		err = db.AutoMigrate(&model.User{}, &model.Blog{}, &model.Portfolio{}, &model.Session{})
		if err != nil {
			log.Fatalf("AutoMigration failed: %v", err)
		}
		log.Println("Auto-migrations completed successfully")
	}

	dbInstance = &service{
		db: db,
	}
	return dbInstance
}

// Health checks the health of the database connection by pinging the database.
func (s *service) Health() map[string]string {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	stats := make(map[string]string)

	sqlDB, err := s.db.DB()
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("failed to get sql.DB: %v", err)
		return stats
	}

	// Ping the database
	err = sqlDB.PingContext(ctx)
	if err != nil {
		stats["status"] = "down"
		stats["error"] = fmt.Sprintf("db down: %v", err)
		log.Fatalf("db down: %v", err)
		return stats
	}

	// Database is up, add stats
	stats["status"] = "up"
	stats["message"] = "It's healthy"

	dbStats := sqlDB.Stats()
	stats["open_connections"] = strconv.Itoa(dbStats.OpenConnections)
	stats["in_use"] = strconv.Itoa(dbStats.InUse)
	stats["idle"] = strconv.Itoa(dbStats.Idle)
	stats["wait_count"] = strconv.FormatInt(dbStats.WaitCount, 10)
	stats["wait_duration"] = dbStats.WaitDuration.String()
	stats["max_idle_closed"] = strconv.FormatInt(dbStats.MaxIdleClosed, 10)
	stats["max_lifetime_closed"] = strconv.FormatInt(dbStats.MaxLifetimeClosed, 10)

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}
	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}
	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}
	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime."
	}

	return stats
}

// Close closes the database connection.
func (s *service) Close() error {
	log.Println("Disconnected from database")
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// GetDB returns the GORM DB instance.
func (s *service) GetDB() *gorm.DB {
	return s.db
}

package database

import (
	"context"
	"fmt"
	"log"
	"log/slog"
	"strconv"
	"time"

	"go-cms/internal/config"
	"go-cms/internal/model"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	gormlogger "gorm.io/gorm/logger"
	"gorm.io/plugin/opentelemetry/tracing"
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
	dbInstance *service
)

func New(cfg *config.Config) Service {
	// Reuse Connection
	if dbInstance != nil {
		return dbInstance
	}

	// Determine GORM log level based on APP_ENV
	gormLogLevel := gormlogger.Silent
	if cfg.AppEnv == "development" {
		gormLogLevel = gormlogger.Info
	}

	// Slow query threshold: queries exceeding this duration will be logged as WARN
	slowQueryThreshold := 200 * time.Millisecond

	gormCfg := &gorm.Config{
		// Replace the default GORM logger with the slog bridge:
		// → structured query logs (JSON) + trace_id + sent to Loki
		Logger: newGormLogger(gormLogLevel, slowQueryThreshold),
	}

	// Retry loop: attempt to connect up to maxDBRetries times.
	// Useful when the app starts before the DB container is fully ready.
	const (
		maxDBRetries = 10
		retryDelay   = 3 * time.Second
	)

	var (
		db  *gorm.DB
		err error
	)
	for attempt := 0; attempt < maxDBRetries; attempt++ {
		db, err = gorm.Open(postgres.Open(cfg.DatabaseURL), gormCfg)
		if err == nil {
			break
		}

		if attempt == 0 {
			log.Printf("CONNECT failed (%v): %v", err, cfg.DatabaseURL)
		} else {
			log.Printf("RECONNECT(%d) failed (%v): %v", attempt, err, cfg.DatabaseURL)
		}

		if attempt < maxDBRetries-1 {
			log.Printf("Retrying in %s... (%d/%d)", retryDelay, attempt+1, maxDBRetries)
			time.Sleep(retryDelay)
		}
	}
	if err != nil {
		log.Fatalf("Failed to connect to database after %d attempts: %v", maxDBRetries, err)
	}

	// Register GORM OTel plugin:
	// → each query automatically becomes a child span in Grafana Tempo
	// → span includes: table name, SQL operation, duration, error (if any)
	if err := db.Use(tracing.NewPlugin(
		tracing.WithoutMetrics(),
	)); err != nil {
		slog.Warn("failed to register gorm otel tracing plugin", "error", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to retrieve generic database object: %v", err)
	}

	// Connection Pooling Settings — sourced from config (reads DB_MAX_IDLE_CONNS / DB_MAX_OPEN_CONNS env).
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	log.Println("Database connection established successfully with GORM")

	// Run GORM AutoMigrate if configured (Optional, usually for dev)
	if cfg.AppEnv == "development" && cfg.AppAutoMigrate {
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

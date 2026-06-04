package config

import (
	"log"
	"os"
	"strconv"

	_ "github.com/joho/godotenv/autoload"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppEnv        string
	AppPort       string
	SecretKey     string
	AppAutoMigrate bool // run GORM AutoMigrate on startup (development only)

	DatabaseURL    string
	DBMaxIdleConns int // default: 10
	DBMaxOpenConns int // default: 100

	RedisURL string

	UploadStorage string // "local" or "s3"
	UploadDir     string
	S3Bucket      string
	S3Region      string
	S3Endpoint    string

	AdminUsername string
	AdminPassword string

	// Telemetry / Observability
	ServiceName    string // Service name displayed in Grafana (default: "go-cms")
	ServiceVersion string // Service version, ideally from git tag (default: "dev")
	OTLPEndpoint   string // Grafana Alloy address: "host:4317" (default: "localhost:4317")
	LogLevel       string // "debug" | "info" | "warn" | "error" (default: "info")
}

// Load reads environment variables and returns a populated Config struct.
func Load() *Config {
	secretKey := os.Getenv("APP_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("APP_SECRET_KEY is required (min 32 chars)")
	}

	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "8080"),
		SecretKey:      secretKey,
		AppAutoMigrate: os.Getenv("APP_AUTO_MIGRATE") == "true",

		DatabaseURL:    requireEnv("DATABASE_URL"),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 100),

		RedisURL: os.Getenv("REDIS_URL"), // optional

		UploadStorage: getEnv("UPLOAD_STORAGE", "local"),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		S3Bucket:      os.Getenv("S3_BUCKET"),
		S3Region:      os.Getenv("S3_REGION"),
		S3Endpoint:    os.Getenv("S3_ENDPOINT"),

		AdminUsername: os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),

		// Telemetry
		ServiceName:    getEnv("SERVICE_NAME", "go-cms"),
		ServiceVersion: getEnv("SERVICE_VERSION", "dev"),
		OTLPEndpoint:   getEnv("OTLP_ENDPOINT", "localhost:4317"),
		LogLevel:       getEnv("LOG_LEVEL", "info"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func requireEnv(key string) string {
	v := os.Getenv(key)
	if v == "" {
		log.Fatalf("required environment variable %q is not set", key)
	}
	return v
}

func getEnvInt(key string, fallback int) int {
	if v, err := strconv.Atoi(os.Getenv(key)); err == nil && v > 0 {
		return v
	}
	return fallback
}

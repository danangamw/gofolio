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

	UploadStorage      string // "local" or "s3"
	UploadDir          string
	S3Bucket           string
	S3Region           string
	S3Endpoint         string
	S3PublicEndpoint   string
	S3AccessKeyID      string
	S3SecretAccessKey  string

	AdminUsername string
	AdminPassword string

	// Rate Limiting
	RateLimitRPS   float64
	RateLimitBurst int

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

	s3AccessKey := getEnv("S3_ACCESS_KEY_ID", os.Getenv("AWS_ACCESS_KEY_ID"))
	s3SecretKey := getEnv("S3_SECRET_ACCESS_KEY", os.Getenv("AWS_SECRET_ACCESS_KEY"))

	rpsStr := getEnv("RATE_LIMIT_RPS", "5")
	rateLimitRPS, err := strconv.ParseFloat(rpsStr, 64)
	if err != nil || rateLimitRPS <= 0 {
		rateLimitRPS = 5.0
	}
	rateLimitBurst := getEnvInt("RATE_LIMIT_BURST", 10)

	return &Config{
		AppEnv:         getEnv("APP_ENV", "development"),
		AppPort:        getEnv("APP_PORT", "8080"),
		SecretKey:      secretKey,
		AppAutoMigrate: os.Getenv("APP_AUTO_MIGRATE") == "true",

		DatabaseURL:    requireEnv("DATABASE_URL"),
		DBMaxIdleConns: getEnvInt("DB_MAX_IDLE_CONNS", 10),
		DBMaxOpenConns: getEnvInt("DB_MAX_OPEN_CONNS", 100),

		RedisURL: os.Getenv("REDIS_URL"), // optional

		UploadStorage:     getEnv("UPLOAD_STORAGE", "local"),
		UploadDir:         getEnv("UPLOAD_DIR", "./uploads"),
		S3Bucket:          getEnv("S3_BUCKET", "go-cms"),
		S3Region:          getEnv("S3_REGION", "us-east-1"),
		S3Endpoint:        os.Getenv("S3_ENDPOINT"),
		S3PublicEndpoint:  os.Getenv("S3_PUBLIC_ENDPOINT"),
		S3AccessKeyID:     s3AccessKey,
		S3SecretAccessKey: s3SecretKey,

		AdminUsername: os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),

		RateLimitRPS:   rateLimitRPS,
		RateLimitBurst: rateLimitBurst,

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

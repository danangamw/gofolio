package config

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	AppEnv    string
	AppPort   string
	SecretKey string

	DatabaseURL string

	RedisURL string

	UploadStorage string // "local" or "s3"
	UploadDir     string
	S3Bucket      string
	S3Region      string
	S3Endpoint    string

	AdminUsername string
	AdminPassword string
}

// Load reads environment variables and returns a populated Config struct.
func Load() *Config {
	secretKey := os.Getenv("APP_SECRET_KEY")
	if secretKey == "" {
		log.Fatal("APP_SECRET_KEY is required (min 32 chars)")
	}

	return &Config{
		AppEnv:    getEnv("APP_ENV", "development"),
		AppPort:   getEnv("APP_PORT", "8080"),
		SecretKey: secretKey,

		DatabaseURL: requireEnv("DATABASE_URL"),

		RedisURL: os.Getenv("REDIS_URL"), // optional

		UploadStorage: getEnv("UPLOAD_STORAGE", "local"),
		UploadDir:     getEnv("UPLOAD_DIR", "./uploads"),
		S3Bucket:      os.Getenv("S3_BUCKET"),
		S3Region:      os.Getenv("S3_REGION"),
		S3Endpoint:    os.Getenv("S3_ENDPOINT"),

		AdminUsername: os.Getenv("ADMIN_USERNAME"),
		AdminPassword: os.Getenv("ADMIN_PASSWORD"),
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

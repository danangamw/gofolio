package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"go-cms/internal/model"
)

func main() {
	dsn := os.Getenv("DATABASE_URL")
	if dsn == "" {
		log.Fatal("DATABASE_URL is required")
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	log.Println("Running migrations...")
	err = db.AutoMigrate(&model.User{}, &model.Blog{}, &model.Portfolio{}, &model.Session{}, &model.SysConfig{})
	if err != nil {
		log.Fatalf("migration failed: %v", err)
	}
	log.Println("Migrations completed successfully.")
}

package main

import (
	"log"
	"os"

	_ "github.com/joho/godotenv/autoload"
)

func main() {
	adminUser := os.Getenv("ADMIN_USERNAME")
	adminPass := os.Getenv("ADMIN_PASSWORD")

	if adminUser == "" || adminPass == "" {
		log.Fatal("ADMIN_USERNAME and ADMIN_PASSWORD are required")
	}

	// TODO: hash password with Argon2id and insert admin user via UserRepository
	log.Printf("Seed: admin user '%s' would be created here.", adminUser)
}

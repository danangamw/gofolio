package main

import (
	"context"
	"log"

	"go-cms/internal/config"
	"go-cms/internal/database"
	"go-cms/internal/model"
	"go-cms/internal/repository"
	"go-cms/internal/service"
)

func main() {
	cfg := config.Load()

	db := database.New(cfg)
	defer db.Close()

	userRepo := repository.NewUserRepository(db.GetDB())

	username := cfg.AdminUsername
	password := cfg.AdminPassword

	if username == "" || password == "" {
		log.Fatal("ADMIN_USERNAME and ADMIN_PASSWORD must be set in .env")
	}

	ctx := context.Background()

	// Check if user already exists — idempotent.
	existing, err := userRepo.FindByUsername(ctx, username)
	if err != nil {
		log.Fatalf("seed: check existing user: %v", err)
	}
	if existing != nil {
		log.Printf("Admin user %q already exists — skipping seed.", username)
		return
	}

	hash, err := service.HashPassword(password)
	if err != nil {
		log.Fatalf("seed: hash password: %v", err)
	}

	user := &model.User{
		Username:     username,
		PasswordHash: hash,
	}

	if err := userRepo.Create(ctx, user); err != nil {
		log.Fatalf("seed: create admin user: %v", err)
	}

	log.Printf("Admin user %q created successfully (id: %s)", username, user.ID)
}

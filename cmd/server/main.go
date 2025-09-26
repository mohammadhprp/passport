package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"

	"github.com/mohammadhprp/passport/internal/config"
	"github.com/mohammadhprp/passport/internal/models"
	"github.com/mohammadhprp/passport/internal/routers"
)

func main() {
	loadDotEnv()

	cfg := config.Load()

	db, err := config.NewPostgresConnection(cfg)
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	if err := models.AutoMigrate(db); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	router := routers.NewRouter(db)

	address := fmt.Sprintf(":%s", cfg.AppPort)
	log.Printf("listening on %s", address)
	if err := router.Run(address); err != nil {
		log.Fatalf("server stopped: %v", err)
	}
}

func loadDotEnv() {
	if _, err := os.Stat(".env"); err == nil {
		if err := godotenv.Load(); err != nil {
			log.Printf("warning: could not load .env: %v", err)
		}
	}
}

package main

import (
	"database/sql"
	api "gemsvietnambe/internal"
	"gemsvietnambe/internal/config"
	"gemsvietnambe/internal/model/database"
	"gemsvietnambe/pkg/logger"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const version = "1.0.0"

func main() {
	err := logger.Init("info", "stdout")
	if err != nil {
		log.Printf("Failed to initialize logger: %v\n", err)
		return
	}

	err = godotenv.Load()
	if err != nil {
		logger.Error("Can't load config from .env. Problem with .env, or the server is in production environment", err)
		return
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Error("Could not connect to db", err)
		return
	}

	dbQueries := database.New(db)

	config := config.ApiConfig{
		Env:     strings.ToUpper(os.Getenv("ENV")),
		Port:    os.Getenv("PORT"),
		Version: os.Getenv("VERSION"),
		DB:      dbQueries,
	}

	server := api.Application{}
	server.Run(config)
}

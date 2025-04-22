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

	secretkey := os.Getenv("SECRETKEY")
	if secretkey == "" {
		logger.Error("SECRETKEY is not set in .env file", nil)
		os.Exit(1)
	}

	refreshkey := os.Getenv("REFRESHKEY")
	if refreshkey == "" {
		logger.Error("REFRESHKEY is not set in .env file", nil)
		os.Exit(1)
	}

	version := os.Getenv("VERSION")
	if version == "" {
		logger.Error("VERSION is not set in .env file", nil)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		logger.Error("DB_URL is not set in .env file", nil)
		os.Exit(1)
	}

	port := os.Getenv("PORT")
	if port == "" {
		logger.Error("PORT is not set in .env file", nil)
		os.Exit(1)
	}

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		logger.Error("Could not connect to db", err)
		os.Exit(1)
		return
	}

	dbQueries := database.New(db)

	config := config.ApiConfig{
		Env:        strings.ToUpper(os.Getenv("ENV")),
		Port:       port,
		Version:    version,
		SecretKey:  secretkey,
		RefreshKey: refreshkey,
		DB:         dbQueries,
	}

	server := api.Application{}
	server.Run(config)
}

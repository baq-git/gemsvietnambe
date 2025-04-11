package api

import (
	"gemsvietnambe/internal/config"
	"gemsvietnambe/internal/handlers"
	"gemsvietnambe/pkg/logger"
	"log"
	"net/http"
	"time"

	"github.com/rs/cors"
)

type Application struct {
	Env     string
	Port    string
	Version string
	handlers.Handlers
}

func (app *Application) Run(apiConfig config.ApiConfig) {
	app.Env = apiConfig.Env
	app.Port = apiConfig.Port
	app.Version = apiConfig.Version
	app.DB = apiConfig.DB

	addr := ":" + app.Port
	if app.Env == "DEV" {
		addr = "localhost:" + app.Port
	}

	handler := cors.Default().Handler(app.routes())

	srv := http.Server{
		Addr:         addr,
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 90 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	err := srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

	logger.Info("Server is listening...")
	srv.ListenAndServe()
}

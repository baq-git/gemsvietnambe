package api

import (
	"gemsvietnambe/internal/config"
	"gemsvietnambe/internal/handlers"
	"gemsvietnambe/internal/middleware"
	"log"
	"net/http"
	"time"

	"github.com/justinas/alice"
	"github.com/rs/cors"
)

type Application struct {
	Env        string
	Port       string
	Version    string
	SecretKey  string
	RefreshKey string
	handlers.Handlers
}

func (app *Application) Run(apiConfig config.ApiConfig) {
	app.Env = apiConfig.Env
	app.Port = apiConfig.Port
	app.Version = apiConfig.Version
	app.SecretKey = apiConfig.SecretKey
	app.RefreshKey = apiConfig.RefreshKey
	app.DB = apiConfig.DB

	addr := ":" + app.Port
	if app.Env == "DEV" {
		addr = "localhost:" + app.Port
	}

	ctxMap := make(map[string]string)
	ctxMap[string(middleware.SecretkeyContextK)] = app.SecretKey
	ctxMap[string(middleware.RefreshkeyContextK)] = app.RefreshKey
	chain := alice.New(middleware.AddContext(ctxMap))

	cors := cors.Default().Handler(app.routes())
	handler := chain.Then(cors)

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
}

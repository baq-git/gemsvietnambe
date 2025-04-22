package api

import (
	"gemsvietnambe/internal/middleware"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	authChain := alice.New(middleware.Authenticate)

	router.Handler(http.MethodGet, "/v1/healthcheck", authChain.ThenFunc(app.HealthCheck))
	router.HandlerFunc(http.MethodPost, "/v1/auth/signup", app.HandlerUserCreate)
	router.HandlerFunc(http.MethodPost, "/v1/auth/login", app.HandlerUserLogin)
	router.HandlerFunc(http.MethodPost, "/v1/auth/refresh", app.HandlerRefreshToken)
	router.Handler(http.MethodPost, "/v1/auth/logout", authChain.ThenFunc(app.HandlerLogout))

	router.HandlerFunc(http.MethodGet, "/v1/gem_categories", app.HandlerGemCategoriesRetrieve)

	router.HandlerFunc(http.MethodPost, "/v1/gems", app.HandlerGemCreate)
	router.Handler(http.MethodGet, "/v1/gems", authChain.ThenFunc(app.HandlerGemsRetrieve))
	router.HandlerFunc(http.MethodGet, "/v1/gems/:id", app.HandlerGetGem)
	router.HandlerFunc(http.MethodDelete, "/v1/gems/:id", app.HandlerGemDelete)
	router.HandlerFunc(http.MethodPatch, "/v1/gems/:id", app.HandlerGemUpdate)

	return router
}

package api

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.HealthCheck)
	router.HandlerFunc(http.MethodPost, "/v1/signup", app.HandlerUserSignUp)
	router.HandlerFunc(http.MethodPost, "/v1/login", app.HandlerUserLogin)

	router.HandlerFunc(http.MethodGet, "/v1/gem_categories", app.HandlerGemCategoriesRetrieve)

	router.HandlerFunc(http.MethodPost, "/v1/gems", app.HandlerGemCreate)
	router.HandlerFunc(http.MethodGet, "/v1/gems", app.HandlerGemsRetrieve)
	router.HandlerFunc(http.MethodDelete, "/v1/gems/:id", app.HandlerGemDelete)
	router.HandlerFunc(http.MethodPatch, "/v1/gems/:id", app.HandlerGemUpdate)

	return router
}

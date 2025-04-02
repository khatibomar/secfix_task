package main

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/latest_data", app.latestDataHandler)
	router.Handler(http.MethodGet, "/debug/vars", expvar.Handler())

	return app.logRequest(app.recoverPanic(app.enableCORS(router)))
}

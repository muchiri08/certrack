package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.index)

	// serving static files
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))

	return router
}

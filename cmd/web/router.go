package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.index)
	router.HandlerFunc(http.MethodGet, "/signup", app.signupForm)
	router.HandlerFunc(http.MethodGet, "/signin", app.signinForm)

	// serving static files
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))

	return router
}

package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/", app.index)
	router.HandlerFunc(http.MethodGet, "/signup", app.signupForm)
	router.HandlerFunc(http.MethodPost, "/signup", app.signup)
	router.HandlerFunc(http.MethodGet, "/signin", app.signinForm)
	router.HandlerFunc(http.MethodPost, "/signin", app.signin)
	router.HandlerFunc(http.MethodGet, "/home/:view", app.home)

	// serving static files
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))

	return router
}

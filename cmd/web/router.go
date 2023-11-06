package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

func (app *application) routes() http.Handler {

	staticMiddleware := alice.New(app.recoverPanic, app.logRequest, app.secureHeaders)
	dynamicMiddleware := alice.New(app.sessionManager.LoadAndSave, app.authenticate)

	router := httprouter.New()

	router.Handler(http.MethodGet, "/", dynamicMiddleware.ThenFunc(app.index))
	router.Handler(http.MethodGet, "/signup", dynamicMiddleware.ThenFunc(app.signupForm))
	router.Handler(http.MethodPost, "/signup", dynamicMiddleware.ThenFunc(app.signup))
	router.Handler(http.MethodGet, "/signin", dynamicMiddleware.ThenFunc(app.signinForm))
	router.Handler(http.MethodPost, "/signin", dynamicMiddleware.ThenFunc(app.signin))
	router.Handler(http.MethodGet, "/home/:view", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.home))
	router.Handler(http.MethodGet, "/logout", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.logout))

	// serving static files
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))

	return staticMiddleware.Then(router)
}

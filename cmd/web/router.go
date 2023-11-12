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
	router.Handler(http.MethodGet, "/logout", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.logout))
	router.Handler(http.MethodGet, "/home/domains", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.getCerts))
	router.Handler(http.MethodGet, "/home/new", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.newDomainForm))
	router.Handler(http.MethodPost, "/home/new", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.newDomain))
	router.Handler(http.MethodGet, "/home/account", dynamicMiddleware.Append(app.requiredAuthentication).ThenFunc(app.account))

	// serving static files
	router.ServeFiles("/static/*filepath", http.Dir("./ui/static"))

	return staticMiddleware.Then(router)
}

package main

import (
	"net/http"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "index.page.html", nil)
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signup.page.html", nil)
}

func (app *application) signinForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "signin.page.html", nil)
}

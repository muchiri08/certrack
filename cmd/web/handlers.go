package main

import (
	"net/http"

	"github.com/muchiri08/certrack/internal/forms"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pINDEX, nil)
}

func (app *application) signupForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pSIGNUP, &templateData{Form: &forms.Form{}})
}

func (app *application) signup(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("email", "username", "password")
	form.MinLength("password", 8)
	form.MatchPattern("email", forms.EmailRegex)
	if !form.Valid() {
		app.render(w, r, pSIGNUP, &templateData{Form: form})
		return
	}

	app.infoLog.Println(form.Get("email"), form.Get("username"), form.Get("password"))

}

func (app *application) signinForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pSIGNIN, &templateData{Form: &forms.Form{}})
}

func (app *application) signin(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}
	form := forms.New(r.PostForm)
	form.Required("username", "password")
	if !form.Valid() {
		app.render(w, r, pSIGNIN, &templateData{Form: form})
		return
	}
	username := form.Get("username")
	password := form.Get("password")

	if username != "kennedy" {
		form.Errors.Add("generic", "Invalid credentials")
		app.render(w, r, pSIGNIN, &templateData{Form: form})
		return
	}

	if password != "12345678" {
		form.Errors.Add("generic", "Invalid credentials")
		app.render(w, r, pSIGNIN, &templateData{Form: form})
		return
	}

}

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	switch {
	case app.getParam(r, "view") == "domains":
		app.render(w, r, pDOMAINS, &templateData{HasSidebar: true})
	case app.getParam(r, "view") == "new":
		app.render(w, r, pNEW, &templateData{HasSidebar: true})
	case app.getParam(r, "view") == "account":
		app.render(w, r, pACCOUNT, &templateData{HasSidebar: true})
	default:
		app.notFound(w)
	}
}

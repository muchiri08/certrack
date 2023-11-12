package main

import (
	"errors"
	"net/http"
	"strings"

	"github.com/muchiri08/certrack/internal/forms"
	"github.com/muchiri08/certrack/internal/models"
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

	user := models.User{
		Username: form.Get("username"),
		Email:    form.Get("email"),
	}

	err := user.Password.Set(form.Get("password"))
	if err != nil {
		app.serverError(w, err)
		return
	}

	err = app.models.Users.NewUser(&user)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrDuplicateEmail):
			form.Errors.Add("generic", "email already exists")
			app.render(w, r, pSIGNUP, &templateData{Form: form})
		case errors.Is(err, models.ErrDuplicateUsername):
			form.Errors.Add("generic", "username already exists")
			app.render(w, r, pSIGNUP, &templateData{Form: form})
		default:
			app.serverError(w, err)
		}
		return
	}

	app.sessionManager.Put(r.Context(), "flash", "Signup successfull. Please signin in below.")

	http.Redirect(w, r, "/signin", http.StatusSeeOther)

}

func (app *application) signinForm(w http.ResponseWriter, r *http.Request) {
	flash := app.sessionManager.PopString(r.Context(), "flash")
	app.render(w, r, pSIGNIN, &templateData{Form: &forms.Form{}, Flash: flash})
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

	user, err := app.models.Users.GetUserByUsername(username)
	if err != nil {
		form.Errors.Add("generic", "Invalid credentials")
		app.render(w, r, pSIGNIN, &templateData{Form: form})
		return
	}

	match, err := user.Password.Matches(password)
	if err != nil {
		app.serverError(w, err)
		return
	}

	if !match {
		form.Errors.Add("generic", "Invalid credentials")
		app.render(w, r, pSIGNIN, &templateData{Form: form})
		return
	}

	app.sessionManager.Put(r.Context(), "username", user.Username)

	http.Redirect(w, r, "/home/domains", http.StatusSeeOther)

}

func (app *application) newDomainForm(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pNEW, &templateData{HasSidebar: true, Form: &forms.Form{}})
}

func (app *application) newDomain(w http.ResponseWriter, r *http.Request) {
	if err := r.ParseForm(); err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form := forms.New(r.PostForm)
	form.Required("domain")

	if !form.Valid() {
		app.render(w, r, pNEW, &templateData{HasSidebar: true, Form: form})
		return
	}

	app.infoLog.Println(form.Get("domain"))

	domainsString := form.Get("domain")
	domains := strings.Split(strings.ReplaceAll(domainsString, " ", ""), ",")

	userId := app.authenticatedUser(r).Id

	certs, err := app.track(userId, domains...)
	if err != nil {
		form.Errors.Add("domain", "One or more domains are invalid.")
		app.render(w, r, pNEW, &templateData{HasSidebar: true, Form: form})
		return
	}

	for _, cert := range certs {
		err := app.models.Certs.Insert(cert)
		if err != nil {
			//TODO: Notify the user what went wrong
			app.serverError(w, err)
			return
		}
	}

	http.Redirect(w, r, "/home/domains", http.StatusSeeOther)
}

func (app *application) account(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, pACCOUNT, &templateData{HasSidebar: true})
}

func (app *application) getCerts(w http.ResponseWriter, r *http.Request) {
	userId := app.authenticatedUser(r).Id
	certs, err := app.models.Certs.GetCerts(userId)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrNoRecordFound):
			app.render(w, r, pDOMAINS, &templateData{HasSidebar: true, Certificates: []*models.Certificate{}})
			return
		default:
			app.serverError(w, err)
			return
		}
	}

	app.render(w, r, pDOMAINS, &templateData{HasSidebar: true, Certificates: certs})
}

func (app *application) logout(w http.ResponseWriter, r *http.Request) {
	app.sessionManager.Remove(r.Context(), "username")
	http.Redirect(w, r, "/signin", http.StatusSeeOther)
}

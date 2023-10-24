package main

import (
	"net/http"
)

func (app *application) index(w http.ResponseWriter, r *http.Request) {
	app.render(w, r, "index.page.html", nil)
}

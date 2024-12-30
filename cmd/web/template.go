package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	ui "github.com/muchiri08/certrack"
	"github.com/muchiri08/certrack/internal/forms"
	"github.com/muchiri08/certrack/internal/models"
)

type templateData struct {
	CurrentYear   int
	HasSidebar    bool
	Form          *forms.Form
	Flash         string
	Authenticated *models.User
	Certificates  []*models.Certificate
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("01 Jan 2006 at 15:04:05")
}

var functions = template.FuncMap{
	"humanDate": humanDate,
}

// caching templates to speed up rendering
func newTemplateCache() (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}
	pages, err := fs.Glob(ui.UI, "ui/html/*.page.html")
	if err != nil {
		return nil, err
	}
	// pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	// if err != nil {
	// 	return nil, err
	// }

	for _, page := range pages {
		// get the name of the template
		name := filepath.Base(page)

		// ts, err := template.ParseFiles(page)
		// if err != nil {
		// 	return nil, err
		// }

		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add partial templates
		ts, err = ts.ParseGlob(filepath.Join("./ui/html", "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// add layout template
		ts, err = ts.ParseGlob(filepath.Join("./ui/html", "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// add the template set to the cache using name of page(like 'index.page.html') as the key.
		cache[name] = ts
	}

	return cache, nil
}

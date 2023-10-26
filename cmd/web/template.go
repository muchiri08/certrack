package main

import (
	"html/template"
	"path/filepath"

	"github.com/muchiri08/certrack/internal/forms"
)

type templateData struct {
	CurrentYear int
	Form        *forms.Form
}

// caching templates to speed up rendering
func newTemplateCache(dir string) (map[string]*template.Template, error) {
	cache := map[string]*template.Template{}

	pages, err := filepath.Glob(filepath.Join(dir, "*.page.html"))
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// get the name of the template
		name := filepath.Base(page)

		ts, err := template.ParseFiles(page)
		if err != nil {
			return nil, err
		}

		// add partial templates
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.partial.html"))
		if err != nil {
			return nil, err
		}

		// add layout template
		ts, err = ts.ParseGlob(filepath.Join(dir, "*.layout.html"))
		if err != nil {
			return nil, err
		}

		// add the template set to the cache using name of page(like 'index.page.html') as the key.
		cache[name] = ts
	}

	return cache, nil
}

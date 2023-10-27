package main

import (
	"html/template"
	"path/filepath"

	"promptbox.tyfacey.net/internal/models"
)

// templateData is struct to pass dynamic data (data from database)
// to HTML templates.
type templateData struct {
	Prompt      models.Prompt
	Prompts     []models.Prompt
	CurrentYear int
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	// Get slice of all filepaths from pages template.
	pages, err := filepath.Glob("./ui/html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Get filename from full path.
		name := filepath.Base(page)

		ts, err := template.ParseFiles("./ui/html/base.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseGlob("./ui/html/partials/*.html")
		if err != nil {
			return nil, err
		}

		ts, err = ts.ParseFiles(page)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

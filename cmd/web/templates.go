package main

import (
	"html/template"
	"path/filepath"

	"promptbox.tyfacey.net/internal/models"
)

type templateData struct {
	Prompt  models.Prompt
	Prompts []models.Prompt
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

		files := []string{
			"./ui/html/base.html",
			"./ui/html/partials/nav.html",
			page,
		}

		// Parse files into template set.
		ts, err := template.ParseFiles(files...)
		if err != nil {
			return nil, err
		}
		cache[name] = ts
	}
	return cache, nil
}

package main

import (
	"html/template"
	"path/filepath"
	"time"

	"promptbox.tyfacey.net/internal/models"
)

// templateData is struct to pass dynamic data (data from database)
// to HTML templates.
type templateData struct {
	Prompt      models.Prompt
	Prompts     []models.Prompt
	CurrentYear int
	Form        any
	Toast       string
}

func humanDate(t time.Time) string {
	return t.Format("Dec 06 2006 at 15:04")
}

// Pass in functions for template use.
var functions = template.FuncMap{
	"humanDate": humanDate,
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

		ts, err := template.New(name).Funcs(functions).ParseFiles("./ui/html/base.html")
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

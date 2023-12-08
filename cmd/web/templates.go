package main

import (
	"html/template"
	"io/fs"
	"path/filepath"
	"time"

	"promptbox.tyfacey.net/internal/models"
	"promptbox.tyfacey.net/ui"
)

// templateData is struct to pass dynamic data (data from database)
// to HTML templates.
type templateData struct {
	Prompt          models.Prompt
	Prompts         []models.Prompt
	CurrentYear     int
	Form            any
	Toast           string
	IsAuthenticated bool
	CSRFToken       string
}

func humanDate(t time.Time) string {
	if t.IsZero() {
		return ""
	}

	return t.UTC().Format("02 Jan 2006 at 15:04")

}

// Pass in functions for template use.
var functions = template.FuncMap{
	"humanDate": humanDate,
}

func newTemplateCache() (map[string]*template.Template, error) {

	cache := map[string]*template.Template{}

	// Get slice of all filepaths from pages template.
	pages, err := fs.Glob(ui.Files, "html/pages/*.html")
	if err != nil {
		return nil, err
	}

	for _, page := range pages {
		// Get filename from full path.
		name := filepath.Base(page)

		patterns := []string{
			"html/base.html",
			"html/partials/*.html",
			page,
		}

		ts, err := template.New(name).Funcs(functions).ParseFS(ui.Files, patterns...)
		if err != nil {
			return nil, err
		}

		cache[name] = ts
	}
	return cache, nil
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"promptbox.tyfacey.net/internal/models"
)

func (app *application) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		app.notFound(w)
		return
	}

	prompts, err := app.prompts.Latest()
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	data := app.newTemplateData(r)
	data.Prompts = prompts

	app.render(w, r, http.StatusOK, "home.html", data)
}

func (app *application) promptView(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id < 1 {
		http.NotFound(w, r)
		return
	}

	prompt, err := app.prompts.Get(id)
	if err != nil {
		if errors.Is(err, models.ErrNoRecord) {
			app.notFound(w)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	data := app.newTemplateData(r)
	data.Prompt = prompt

	app.render(w, r, http.StatusOK, "view.html", data)
}

func (app *application) promptCreate(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}

	// Variables to test prompt creation
	title := "Revere a linked list"
	content := "Write a code snippet to revere a linked list"
	expires := 7

	id, err := app.prompts.Insert(title, content, expires)

	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/prompt/view?id=%d", id), http.StatusSeeOther)
}

package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"unicode/utf8"

	"github.com/julienschmidt/httprouter"
	"promptbox.tyfacey.net/internal/models"
)

type promptCreateForm struct {
	Title       string
	Content     string
	Expires     int
	FieldErrors map[string]string
}

func (app *application) home(w http.ResponseWriter, r *http.Request) {

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

	params := httprouter.ParamsFromContext(r.Context())

	id, err := strconv.Atoi(params.ByName("id"))
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

func (app *application) promptCreatePost(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	expires, err := strconv.Atoi(r.PostForm.Get("expires"))
	if err != nil {
		// Send a 400 response if string to number conversion fails.
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Get title, content, and expiry from request body.
	form := promptCreateForm{
		Title:       r.PostForm.Get("title"),
		Content:     r.PostForm.Get("content"),
		Expires:     expires,
		FieldErrors: map[string]string{}, // empty map to hold any validatioon errors.
	}

	// Check title field
	if strings.TrimSpace(form.Title) == "" {
		form.FieldErrors["title"] = "This cannot be blank"
	} else if utf8.RuneCountInString(form.Title) > 100 {
		form.FieldErrors["title"] = "Your title cannot be more than 100 characters long"
	}
	// Check content field
	if strings.TrimSpace(form.Content) == "" {
		form.FieldErrors["content"] = "This cannot be blank"
	} else if utf8.RuneCountInString(form.Content) > 3000 {
		form.FieldErrors["content"] = "Your prompt cannot be more than 3000 characters long"
	}

	// Check expires field
	if form.Expires != 1 && form.Expires != 7 && form.Expires != 365 {
		form.FieldErrors["expires"] = "You must put in an expire date for your prompt"
	}

	if len(form.FieldErrors) > 0 {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "create.html", data)
		return
	}

	id, err := app.prompts.Insert(form.Title, form.Content, form.Expires)
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	http.Redirect(w, r, fmt.Sprintf("/prompt/view/%d", id), http.StatusSeeOther)
}

func (app *application) promptCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = promptCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.html", data)
}

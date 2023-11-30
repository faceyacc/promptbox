package main

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
	"promptbox.tyfacey.net/internal/models"
	"promptbox.tyfacey.net/internal/validator"
)

type promptCreateForm struct {
	Title               string `form:"title"`
	Content             string `form:"content"`
	Expires             int    `form:"expires"`
	validator.Validator `form:"-"`
}

type userSignUpForm struct {
	Name                string `form:"name"`
	Email               string `form:"email"`
	Password            string `form:"password"`
	validator.Validator `form:"-"`
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

	var form promptCreateForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	// Check title field
	form.CheckField(validator.NotBlank(form.Title), "title", "Title cannot be blank")
	form.CheckField(validator.MaxChars(form.Title, 100), "title", "Title cannot be more than 100 characters long")

	// Check content field
	form.CheckField(validator.NotBlank(form.Content), "content", "Prompt cannot be blank")
	form.CheckField(validator.MaxChars(form.Content, 3000), "content", "Prompts cannot be more than 3000 characters long")

	// Check expires field
	form.CheckField(validator.PermittedValue(form.Expires, 1, 7, 365), "expires", "You must put in an expire date for your prompt")

	if !form.Valid() {
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

	// Add session data.
	app.sessionManager.Put(r.Context(), "toast", "Prompt ssuccesfully created!")

	http.Redirect(w, r, fmt.Sprintf("/prompt/view/%d", id), http.StatusSeeOther)
}

func (app *application) promptCreate(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = promptCreateForm{
		Expires: 365,
	}

	app.render(w, r, http.StatusOK, "create.html", data)
}

// Handlers to handle user authentication.
func (app *application) userSignup(w http.ResponseWriter, r *http.Request) {

	data := app.newTemplateData(r)

	data.Form = userSignUpForm{}

	app.render(w, r, http.StatusOK, "signup.html", data)
}

func (app *application) userSignupPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Create a new user...")
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Display a HTML form for logging in a user")
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Authenticate and login the user...")
}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Logout the user")
}

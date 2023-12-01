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

type userLoginForm struct {
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
	var form userSignUpForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Name), "name", "This field cannot be blank")
	form.CheckField(validator.NotBlank(form.Email), "email", "You must input your email to signup")
	form.CheckField(validator.NotBlank(form.Password), "password", "You must create a password")

	form.CheckField(validator.Matches(form.Email), "email", "You must give an valid email")
	form.CheckField(validator.MinChars(form.Password, 8), "password", "Your password must be at least 8 characters long")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		return
	}

	err = app.users.Insert(form.Name, form.Email, form.Password)
	// Handle if user tries to sign up with duplicate email.
	if err != nil {
		if errors.Is(err, models.ErrDuplicateEmail) {
			form.AddFieldError("email", "This email is already in use")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "signup.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	// Show user success flash message.
	app.sessionManager.Put(r.Context(), "flash", "Your signup was successful. Please log in.")

	// Redirect user to the login page if sign up is successful.
	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

func (app *application) userLogin(w http.ResponseWriter, r *http.Request) {
	data := app.newTemplateData(r)

	data.Form = userLoginForm{}

	app.render(w, r, http.StatusOK, "login.html", data)
}

func (app *application) userLoginPost(w http.ResponseWriter, r *http.Request) {
	var form userLoginForm

	err := app.decodePostForm(r, &form)
	if err != nil {
		app.clientError(w, http.StatusBadRequest)
		return
	}

	form.CheckField(validator.NotBlank(form.Email), "email", "You must input your email to signin")
	form.CheckField(validator.NotBlank(form.Password), "password", "This field cannot be blank")

	form.CheckField(validator.Matches(form.Email), "email", "This field must be a valid email address")

	if !form.Valid() {
		data := app.newTemplateData(r)
		data.Form = form
		app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		return
	}

	id, err := app.users.Authenticate(form.Email, form.Password)
	if err != nil {
		if errors.Is(err, models.ErrInvalidCreds) {
			form.AddNonFieldErrors("Email or Password is incorrect")

			data := app.newTemplateData(r)
			data.Form = form
			app.render(w, r, http.StatusUnprocessableEntity, "login.html", data)
		} else {
			app.serverError(w, r, err)
		}
		return
	}

	err = app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	app.sessionManager.Put(r.Context(), "authenticateUserID", id)

	http.Redirect(w, r, "/prompt/create", http.StatusSeeOther)

}

func (app *application) userLogoutPost(w http.ResponseWriter, r *http.Request) {
	// Renew token to change sessions ID.
	err := app.sessionManager.RenewToken(r.Context())
	if err != nil {
		app.serverError(w, r, err)
		return
	}

	// Remove the authenticatedUserID from session data.
	app.sessionManager.Remove(r.Context(), "authenticatedUserID")

	// Notify user they are logged out.
	app.sessionManager.Put(r.Context(), "flash", "You've been logged out")

	x := app.sessionManager.Get(r.Context(), "flash")
	fmt.Print(x)

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

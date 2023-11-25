package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/justinas/alice"
)

// Return a servemux containing application routes.
func (app *application) routes() http.Handler {

	router := httprouter.New()
	router.NotFound = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		app.notFound(w)
	})

	fileServer := http.FileServer(http.Dir("./ui/static/"))
	router.Handler(http.MethodGet, "/static/*filepath", http.StripPrefix("/static", fileServer))

	dynamic := alice.New(app.sessionManager.LoadAndSave)

	router.Handler(http.MethodGet, "/", dynamic.ThenFunc(app.home))
	router.Handler(http.MethodGet, "/prompt/view/:id", dynamic.ThenFunc(app.promptView))
	router.Handler(http.MethodGet, "/prompt/create", dynamic.ThenFunc(app.promptCreate))
	router.Handler(http.MethodPost, "/prompt/create", dynamic.ThenFunc(app.promptCreatePost))

	// Routes to handle user authentication.
	router.Handler(http.MethodGet, "/user/signup", dynamic.ThenFunc(app.userSignup))
	router.Handler(http.MethodPost, "/user/signup", dynamic.ThenFunc(app.userSignupPost))
	router.Handler(http.MethodGet, "/user/login", dynamic.ThenFunc(app.userLogin))
	router.Handler(http.MethodPost, "/user/login", dynamic.ThenFunc(app.userLoginPost))
	router.Handler(http.MethodPost, "/user/logout", dynamic.ThenFunc(app.userLogoutPost))

	requestMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return requestMiddleware.Then(router)
}

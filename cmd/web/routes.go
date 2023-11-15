package main

import (
	"net/http"

	"github.com/justinas/alice"
)

// Return a servemux containing application routes.
func (app *application) routes() http.Handler {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/prompt/view", app.promptView)
	mux.HandleFunc("/prompt/create", app.promptCreate)

	requestMiddleware := alice.New(app.recoverPanic, app.logRequest, secureHeaders)

	return requestMiddleware.Then(mux)
}

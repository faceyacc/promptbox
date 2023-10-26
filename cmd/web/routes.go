package main

import "net/http"

// Return a servemux containing application routes.
func (app *application) routes() *http.ServeMux {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", app.home)
	mux.HandleFunc("/prompt/view", app.promptView)
	mux.HandleFunc("/prompt/create", app.promptCreate)

	return mux
}

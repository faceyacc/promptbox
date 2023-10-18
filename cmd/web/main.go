package main

import (
	"log"
	"net/http"
)

func main() {
	mux := http.NewServeMux()

	fileServer := http.FileServer(http.Dir("./ui/static/"))

	mux.Handle("/static/", http.StripPrefix("/static/", fileServer))
	mux.HandleFunc("/", home)
	mux.HandleFunc("/prompt/view", promptView)
	mux.HandleFunc("/prompt/create", promptCreate)

	log.Print("Starting seerver on: http://localhost:4000/")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}
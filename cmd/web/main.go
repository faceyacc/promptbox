package main

import (
	"flag"
	"log/slog"
	"net/http"
	"os"
)

// Struct to hold application-wide dependencies.
type application struct {
	logger *slog.Logger
}

func main() {

	addr := flag.String("addr", ":4000", "HTTP network address")
	flag.Parse()

	// Added for structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	app := &application{
		logger: logger,
	}

	logger.Info("starting server", slog.String("addr", *addr))

	err := http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

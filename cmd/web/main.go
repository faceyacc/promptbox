package main

import (
	"database/sql"
	"flag"
	"html/template"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"promptbox.tyfacey.net/internal/models"
)

// Struct to hold application-wide dependencies.
type application struct {
	logger        *slog.Logger
	prompts       *models.PromptModel
	templateCache map[string]*template.Template
}

func main() {
	addr := flag.String("addr", ":4000", "HTTP network address")
	dsn := flag.String("dsn", "web:turintest@/promptbox?parseTime=true", "MySQL data source name")
	flag.Parse()

	// Added for structured logging
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	// Creates a database connection pool
	db, err := openDB(*dsn)
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// Initalize template cache
	templateCache, err := newTemplateCache()
	if err != nil {
		logger.Error(err.Error())
		os.Exit(1)
	}

	// close connection pool before main() exits
	defer db.Close()

	app := &application{
		logger:        logger,
		prompts:       &models.PromptModel{DB: db},
		templateCache: templateCache,
	}

	logger.Info("starting server", "addr", *addr)

	err = http.ListenAndServe(*addr, app.routes())

	logger.Error(err.Error())
	os.Exit(1)
}

func openDB(dsn string) (*sql.DB, error) {
	db, err := sql.Open("mysql", dsn)

	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}

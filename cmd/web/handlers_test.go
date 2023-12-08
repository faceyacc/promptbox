package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"

	"promptbox.tyfacey.net/internal/assert"
)

func TestPing(t *testing.T) {

	app := &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}

	// Test all application handlers.
	server := httptest.NewTLSServer(app.routes())
	defer server.Close()

	// Spin up server client to make request to test server.
	response, err := server.Client().Get(server.URL + "/ping")
	if err != nil {
		t.Fatal(err)
	}

	// Check value of response status code returns a 200.
	assert.Equal(t, response.StatusCode, http.StatusOK)

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	// Check returned string from ping.
	assert.Equal(t, string(body), "OK")

}

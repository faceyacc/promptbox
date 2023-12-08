package main

import (
	"net/http"
	"testing"

	"promptbox.tyfacey.net/internal/assert"
)

func TestPing(t *testing.T) {

	app := newTestApplication(t)

	// Test all application handlers.
	server := newTestServer(t, app.routes())
	defer server.Close()

	// Spin up server client to make request to test server.
	resInt, _, body := server.get(t, "/ping")

	// Check value of response status code returns a 200.
	assert.Equal(t, resInt, http.StatusOK)

	// Check returned string from ping.
	assert.Equal(t, string(body), "OK")

}

package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"promptbox.tyfacey.net/internal/assert"
)

func TestPing(t *testing.T) {
	recorder := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Call test handler.
	ping(recorder, r)

	recorderRes := recorder.Result()

	// Check status code from ping.
	assert.Equal(t, recorderRes.StatusCode, http.StatusOK)

	defer recorderRes.Body.Close()

	body, err := io.ReadAll(recorderRes.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)

	// Check returned string from ping.
	assert.Equal(t, string(body), "OK")

}

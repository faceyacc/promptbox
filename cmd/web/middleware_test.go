package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"promptbox.tyfacey.net/internal/assert"
)

func TestSecureHeaders(t *testing.T) {

	recorder := httptest.NewRecorder()

	r, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	// Mock handler to pass into secureHeaders middleware.
	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	secureHeaders(next).ServeHTTP(recorder, r)

	recorderRes := recorder.Result()

	// Check if middleware has set correct headers based on response.
	expectedVal := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"
	assert.Equal(t, recorderRes.Header.Get("Content-Security-Policy"), expectedVal)

	expectedVal = "origin-when-cross-origin"
	assert.Equal(t, recorderRes.Header.Get("Referrer-Policy"), expectedVal)

	expectedVal = "nosniff"
	assert.Equal(t, recorderRes.Header.Get("X-Content-Type-Options"), expectedVal)

	expectedVal = "deny"
	assert.Equal(t, recorderRes.Header.Get("X-Frame-Options"), expectedVal)

	expectedVal = "0"
	assert.Equal(t, recorderRes.Header.Get("X-XSS-Protection"), expectedVal)

	assert.Equal(t, recorderRes.StatusCode, http.StatusOK)

	defer recorderRes.Body.Close()

	body, err := io.ReadAll(recorderRes.Body)
	if err != nil {
		t.Fatal(err)
	}

	body = bytes.TrimSpace(body)
	assert.Equal(t, string(body), "OK")
}

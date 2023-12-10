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

func TestPromptView(t *testing.T) {
	app := newTestApplication(t)

	testServer := newTestServer(t, app.routes())
	defer testServer.Close()

	tests := []struct {
		name     string
		urlPath  string
		wantCode int
		wantBody string
	}{
		{
			name:     "Valid ID",
			urlPath:  "/prompt/view/1",
			wantCode: http.StatusOK,
			wantBody: "Explain to me like a 5 year-old",
		},
		{
			name:     "Non-existent ID",
			urlPath:  "/prompt/view/2",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Negative ID",
			urlPath:  "/prompt/view/-1",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "String ID",
			urlPath:  "/prompt/view/foo",
			wantCode: http.StatusNotFound,
		},
		{
			name:     "Empty ID",
			urlPath:  "/prompt/view/",
			wantCode: http.StatusNotFound,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			code, _, body := testServer.get(t, test.urlPath)
			assert.Equal(t, code, test.wantCode)

			if test.wantBody != "" {
				assert.StringContains(t, body, test.wantBody)
			}
		})
	}

}

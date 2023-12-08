package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testServer struct {
	*httptest.Server
}

// newTestApplication returns an instance of an application
// stuct with mock dependencies.
func newTestApplication(t *testing.T) *application {
	return &application{
		logger: slog.New(slog.NewTextHandler(io.Discard, nil)),
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	tServer := httptest.NewTLSServer(h)
	return &testServer{tServer}
}

func (tServer *testServer) get(t *testing.T, urlPath string) (int, http.Header, string) {
	response, err := tServer.Client().Get(tServer.URL + urlPath)
	if err != nil {
		t.Fatal(err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	if err != nil {
		t.Fatal(err)
	}
	body = bytes.TrimSpace(body)
	return response.StatusCode, response.Header, string(body)
}

package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"promptbox.tyfacey.net/internal/models/mocks"
)

type testServer struct {
	*httptest.Server
}

// newTestApplication returns an instance of an application
// stuct with mock dependencies.
func newTestApplication(t *testing.T) *application {
	templateCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()

	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	return &application{
		logger:         slog.New(slog.NewTextHandler(io.Discard, nil)),
		templateCache:  templateCache,
		formDecoder:    formDecoder,
		sessionManager: sessionManager,
		prompts:        &mocks.PromptModel{},
		users:          &mocks.UserModel{},
	}
}

func newTestServer(t *testing.T, h http.Handler) *testServer {
	tServer := httptest.NewTLSServer(h)

	// Initalize cookiejar to store any HTTP response cookies to the test server client.
	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}
	tServer.Client().Jar = jar

	tServer.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

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

package api_test

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/zorcal/sbgfit/backend/api"
)

func testServer(t *testing.T, cfg api.Config) *httptest.Server {
	t.Helper()

	h, err := api.NewHandler(cfg)
	if err != nil {
		t.Fatalf("failed to create new handler: %v", err)
	}

	srv := httptest.NewServer(h)

	return srv
}

func makeRequest(t *testing.T, srv *httptest.Server, method, path string, body io.Reader) *http.Response {
	t.Helper()

	req, err := http.NewRequestWithContext(t.Context(), method, srv.URL+path, body)
	if err != nil {
		t.Fatalf("failed to create new request: %v", err)
	}

	resp, err := srv.Client().Do(req)
	if err != nil {
		t.Fatalf("failed to make request: %v", err)
	}
	t.Cleanup(func() { resp.Body.Close() })

	return resp
}

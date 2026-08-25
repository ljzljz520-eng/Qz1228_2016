package nebula

import (
	"gesture-nebula-console/api"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestAPIHandlers(t *testing.T) {
	_, svc := openFixture(t)
	server := api.NewServer(svc, nil)
	request := httptest.NewRequest(http.MethodPost, "/records", strings.NewReader(`{"id":"api1","owner_id":"u1","label":"rose"}`))
	request.Header.Set("Content-Type", "application/json")
	response := httptest.NewRecorder()
	server.Handler().ServeHTTP(response, request)
	if response.Code != http.StatusCreated {
		t.Fatalf("status %d body %s", response.Code, response.Body.String())
	}
	get := httptest.NewRequest(http.MethodGet, "/records/api1", nil)
	getResponse := httptest.NewRecorder()
	server.Handler().ServeHTTP(getResponse, get)
	if getResponse.Code != http.StatusOK {
		t.Fatalf("get status %d", getResponse.Code)
	}
	_, _ = io.Copy(io.Discard, getResponse.Body)
}

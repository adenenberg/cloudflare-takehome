package main

import (
	"cloudflare-takehome/routes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var rt = routes.Routes()

func TestPingEndpoint(t *testing.T) {
	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	rec := httptest.NewRecorder()
	rt.ServeHTTP(rec, req)

	if status := rec.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	if !strings.Contains(rec.Body.String(), "pong") {
		t.Errorf("returned unexpected body: got %v", rec.Body.String())
	}
}

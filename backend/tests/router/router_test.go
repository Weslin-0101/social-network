package router

import (
	"backend/src/router"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGenerate_NotNil(t *testing.T) {
	r := router.Generate()
	if r == nil {
		t.Fatal("expected router to be not nil")
	}
}

func TestGenerate_RoutesExist(t *testing.T) {
	r := router.Generate()

	tests := []struct {
		method string
		path   string
	}{
		{"POST", "/users"},
		{"GET", "/users"},
		{"GET", "/users/1"},
		{"PUT", "/users/1"},
		{"DELETE", "/users/1"},
	}

	for _, paths := range tests {
		req := httptest.NewRequest(paths.method, paths.path, nil)
		rr := httptest.NewRecorder()
		r.ServeHTTP(rr, req)

		if rr.Code == http.StatusNotFound {
			t.Errorf("route %s %s not found (got 404)", paths.method, paths.path)
		}
	}
}

func TestGenerate_MethodNotAllowed(t *testing.T) {
	r := router.Generate()

	req := httptest.NewRequest(http.MethodPatch, "/users/1", nil)
	rr := httptest.NewRecorder()
	r.ServeHTTP(rr, req)

	if rr.Code != http.StatusMethodNotAllowed && rr.Code != http.StatusNotFound {
		t.Errorf("expected 405 or 404 for PATCH /users/1, got %d", rr.Code)
	}
}
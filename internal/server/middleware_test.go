package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHeaderMW(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	sut := HeaderMW(handler)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://test.com", nil)
	sut.ServeHTTP(w, r)

	if w.Header().Get("Access-Control-Allow-Origin") != "*" {
		t.Errorf("expected *, got %s", w.Header().Get("Access-Control-Allow-Origin"))
	}

	if w.Header().Get("Content-Type") != "application/json" {
		t.Errorf("expected *, got %s", w.Header().Get("Content-Type"))
	}
}

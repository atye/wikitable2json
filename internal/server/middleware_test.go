package server

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
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
func TestRequestValidationAndMetricsMW(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mp := &mockPublisher{}
	sut := RequestValidationAndMetricsMW(handler, mp)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://test.com/api/page?verbose=true", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

loop:
	for {
		select {
		case <-ctx.Done():
			t.Errorf("timed out waiting for publish to be called")
		default:
			if !mp.getPublishedCalled() {
				time.Sleep(500 * time.Millisecond)
				continue
			}
			break loop
		}
	}
}

func TestRequestValidationAndMetricsMWEmptyPage(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mp := &mockPublisher{}
	sut := RequestValidationAndMetricsMW(handler, mp)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://test.com/api?verbose=true", nil)
	r.SetPathValue("page", "")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestRequestValidationAndMetricsMWBadParameter(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {})

	mp := &mockPublisher{}
	sut := RequestValidationAndMetricsMW(handler, mp)

	w := httptest.NewRecorder()
	r := httptest.NewRequest(http.MethodGet, "http://test.com/api/page?table=x", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
}

type mockPublisher struct {
	publishCalled bool
	lock          sync.Mutex
}

func (m *mockPublisher) Publish(code int, ip string, page string, lang string) error {
	m.lock.Lock()
	defer m.lock.Unlock()
	m.publishCalled = true
	return nil
}

func (m *mockPublisher) getPublishedCalled() bool {
	m.lock.Lock()
	defer m.lock.Unlock()
	return m.publishCalled
}

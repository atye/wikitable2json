package api

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/atye/wikitable2json/internal/status"
)

var data = []byte(`
<!DOCTYPE html>
<html>
<body>
    <table class="wikitable">
        <tbody>
        </tbody>
    </table>
</body>
</html>
`)

func TestWikiClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/test":
				w.Write([]byte(data))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		originalGetEndpoint := GetEndpoint
		GetEndpoint = func(_, page string) string { return fmt.Sprintf("%s/%s", ts.URL, page) }
		defer func() {
			GetEndpoint = originalGetEndpoint
		}()

		sut := NewWikiClient()

		got, err := sut.GetPage(context.Background(), "test", "en", "")
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}

		want := string(data)

		if want != string(got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("Error from Wikipedia API", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/test":
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				w.Write([]byte("error"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		originalGetEndpoint := GetEndpoint
		GetEndpoint = func(_, page string) string { return fmt.Sprintf("%s/%s", ts.URL, page) }
		defer func() {
			GetEndpoint = originalGetEndpoint
		}()

		sut := NewWikiClient()

		_, got := sut.GetPage(context.Background(), "test", "en", "")

		want := status.Status{
			Message: "error",
			Code:    http.StatusRequestEntityTooLarge,
			Details: status.Details{
				status.Page: "test",
			},
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("WithHTTPClient", func(t *testing.T) {
		sut := NewWikiClient(WithHTTPClient(&http.Client{Timeout: 20 * time.Second}))

		want := 20 * time.Second
		if sut.client.Timeout != want {
			t.Errorf("expected client timeout %v, got %v", want, sut.client.Timeout)
		}
	})
}

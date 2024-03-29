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
			case "/api/rest_v1/page/html/test":
				w.Write([]byte(data))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		sut := NewWikiClient(ts.URL)

		got, err := sut.GetPageBytes(context.Background(), "test", "en", "")
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
			case "/api/rest_v1/page/html/test":
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				w.Write([]byte("error"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		sut := NewWikiClient(ts.URL)

		_, got := sut.GetPageBytes(context.Background(), "test", "en", "")

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

	t.Run("BuildURL", func(t *testing.T) {
		tests := []struct {
			endpoint string
			page     string
			lang     string
			want     string
		}{
			{
				BaseURL,
				"test",
				"en",
				fmt.Sprintf(BaseURL, "en", "test"),
			}, {
				"http://127.0.0.1:61051",
				"test",
				"en",
				"http://127.0.0.1:61051/api/rest_v1/page/html/test",
			},
		}

		for _, tc := range tests {
			got, err := buildURL(tc.endpoint, tc.page, tc.lang)
			if err != nil {
				t.Fatal(err)
			}

			if tc.want != got {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		}
	})

	t.Run("WithHTTPClient", func(t *testing.T) {
		sut := NewWikiClient("", WithHTTPClient(&http.Client{Timeout: 20 * time.Second}))

		want := 20 * time.Second
		if sut.client.Timeout != want {
			t.Errorf("expected client timeout %v, got %v", want, sut.client.Timeout)
		}
	})
}

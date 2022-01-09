package data

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/atye/wikitable-api/internal/status"
)

func TestWikiClient(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/api/rest_v1/page/html/test":
				w.Write([]byte("test"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		sut := NewWikiClient(ts.URL)

		r, err := sut.GetPageData(context.Background(), "test", "en", "")
		if err != nil {
			t.Errorf("expected nil error, got %v", err)
		}

		want := "test"

		got, err := io.ReadAll(r)
		if err != nil {
			t.Fatal(err)
		}

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

		_, got := sut.GetPageData(context.Background(), "test", "en", "")

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
}

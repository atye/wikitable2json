package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/atye/wikitable2json/pkg/client/status"
)

func TestBuildCacheKey(t *testing.T) {
	tests := []struct {
		name string
		qv   queryValues
		page string
		want string
	}{
		{
			"en-0-false-2-false-false",
			queryValues{
				lang:      "en",
				tables:    []int{0},
				cleanRef:  false,
				keyRows:   2,
				verbose:   false,
				brNewLine: false,
			},
			"test",
			"test-en-0-false-2-false-false",
		},
		{
			"en-0-true-2-true-true",
			queryValues{
				lang:      "en",
				tables:    []int{0},
				cleanRef:  true,
				keyRows:   2,
				verbose:   true,
				brNewLine: true,
			},
			"test",
			"test-en-0-true-2-true-true",
		},
		{
			"en-01-true-2-true-true",
			queryValues{
				lang:      "en",
				tables:    []int{0, 1},
				cleanRef:  true,
				keyRows:   2,
				verbose:   true,
				brNewLine: true,
			},
			"test",
			"test-en-01-true-2-true-true",
		},
		{
			"en-all-true-2-true-true",
			queryValues{
				lang:      "en",
				cleanRef:  true,
				keyRows:   2,
				verbose:   true,
				brNewLine: true,
			},
			"test",
			"test-en-all-true-2-true-true",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got := buildCacheKey(tc.page, tc.qv)
			if got != tc.want {
				t.Errorf("want %v, got %v", tc.want, got)
			}
		})
	}
}

func TestParseParameters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		params := r.URL.Query()
		params.Add("lang", "sp")
		params.Add("table", "0")
		params.Add("format", "keyValue")
		params.Add("cleanRef", "true")
		params.Add("keyRows", "2")
		params.Add("verbose", "true")
		r.URL.RawQuery = params.Encode()

		qv, err := parseParameters(r)
		if err != nil {
			t.Fatal(err)
		}

		gotLang := qv.lang
		gotTables := qv.tables
		gotCleanRef := qv.cleanRef
		gotKeyRows := qv.keyRows
		gotVerbose := qv.verbose

		wantLang := "sp"
		wantTables := []int{0}
		wantCleanRef := true
		wantKeyRows := 2
		wantVerbose := true

		if wantLang != gotLang {
			t.Errorf("want %v, got %v", wantLang, gotLang)
		}

		if wantCleanRef != gotCleanRef {
			t.Errorf("want %v, got %v", wantCleanRef, gotCleanRef)
		}

		if !reflect.DeepEqual(wantTables, gotTables) {
			t.Errorf("want %v, got %v", wantTables, gotTables)
		}

		if wantKeyRows != gotKeyRows {
			t.Errorf("want %d, got %d", wantKeyRows, gotKeyRows)
		}

		if wantVerbose != gotVerbose {
			t.Errorf("want %v, got %v", wantVerbose, gotVerbose)
		}
	})

	t.Run("Bad table query", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		params := r.URL.Query()
		params.Add("table", "x")
		r.URL.RawQuery = params.Encode()

		_, got := parseParameters(r)
		if got == nil {
			t.Fatal("expected non-nil error")
		}

		want := status.NewStatus(`strconv.Atoi: parsing "x": invalid syntax`, http.StatusBadRequest)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})

	t.Run("Bad keyrows query", func(t *testing.T) {
		t.Run("Bad keyrows syntax", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api", nil)
			params := r.URL.Query()
			params.Add("keyRows", "x")
			r.URL.RawQuery = params.Encode()

			_, got := parseParameters(r)
			if got == nil {
				t.Fatal("expected non-nil error")
			}

			want := status.NewStatus(`strconv.Atoi: parsing "x": invalid syntax`, http.StatusBadRequest)
			if !reflect.DeepEqual(want, got) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})

		t.Run("keyrows less than 1", func(t *testing.T) {
			r := httptest.NewRequest(http.MethodGet, "/api", nil)
			params := r.URL.Query()
			params.Add("keyRows", "0")
			r.URL.RawQuery = params.Encode()

			_, got := parseParameters(r)
			if got == nil {
				t.Fatal("expected non-nil error")
			}

			want := status.NewStatus(`keyRows must be at least 1`, http.StatusBadRequest)
			if !reflect.DeepEqual(want, got) {
				t.Errorf("expected %v, got %v", want, got)
			}
		})
	})
}

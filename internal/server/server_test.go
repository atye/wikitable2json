package server

import (
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/atye/wikitable-api/internal/status"
)

func TestParseParameters(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		params := r.URL.Query()
		params.Add("lang", "sp")
		params.Add("table", "0")
		params.Add("format", "verbose")
		r.URL.RawQuery = params.Encode()

		gotLang, gotFormat, gotTables, err := parseParameters(r)
		if err != nil {
			t.Fatal(err)
		}

		wantLang := "sp"
		wantFormat := "verbose"
		wantTables := []int{0}

		if wantLang != gotLang {
			t.Errorf("expected %v, got %v", wantLang, gotLang)
		}

		if wantFormat != gotFormat {
			t.Errorf("expected %v, got %v", wantFormat, gotFormat)
		}

		if !reflect.DeepEqual(wantTables, gotTables) {
			t.Errorf("expected %v, got %v", wantTables, gotTables)
		}
	})

	t.Run("Error", func(t *testing.T) {
		r := httptest.NewRequest(http.MethodGet, "/api", nil)
		params := r.URL.Query()
		params.Add("table", "x")
		r.URL.RawQuery = params.Encode()

		_, _, _, got := parseParameters(r)
		if got == nil {
			t.Fatal("expected non-nil error")
		}

		want := status.NewStatus(`strconv.Atoi: parsing "x": invalid syntax`, http.StatusBadRequest)
		if !reflect.DeepEqual(want, got) {
			t.Errorf("expected %v, got %v", want, got)
		}
	})
}

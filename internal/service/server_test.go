package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"
)

func TestServerSuccess(t *testing.T) {
	t.Run("table", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "table.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := []Table{
			{
				Caption: "test",
				Data: [][]string{
					{"Column 1", "Column 2", "Column 3"},
					{"A", "B", "B"},
					{"A", "C", "D"},
					{"E", "F", "F"},
					{"G", "F", "F"},
					{"H", "H", "H"},
				},
			},
			{
				Data: [][]string{
					{"Column 1", "Column 2", "Column 3"},
				},
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/table", nil)

		svr.ServeHTTP(w, r)

		var got []Table
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("first table", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "table.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := []Table{
			{
				Caption: "test",
				Data: [][]string{
					{"Column 1", "Column 2", "Column 3"},
					{"A", "B", "B"},
					{"A", "C", "D"},
					{"E", "F", "F"},
					{"G", "F", "F"},
					{"H", "H", "H"},
				},
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/table", nil)
		params := r.URL.Query()
		params.Add("table", "0")
		r.URL.RawQuery = params.Encode()

		svr.ServeHTTP(w, r)

		var got []Table
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("non-english lang", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "table.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := []Table{
			{
				Caption: "test",
				Data: [][]string{
					{"Column 1", "Column 2", "Column 3"},
					{"A", "B", "B"},
					{"A", "C", "D"},
					{"E", "F", "F"},
					{"G", "F", "F"},
					{"H", "H", "H"},
				},
			},
			{
				Data: [][]string{
					{"Column 1", "Column 2", "Column 3"},
				},
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/table", nil)
		params := r.URL.Query()
		params.Add("lang", "cs")
		r.URL.RawQuery = params.Encode()

		svr.ServeHTTP(w, r)

		var got []Table
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("issue1", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "issue_1.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := []Table{
			{
				Data: [][]string{
					{"Jeju", "South Korea", "official, in Jeju Island"},
					{"Jeju"},
				},
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/issue_1", nil)

		svr.ServeHTTP(w, r)

		var got []Table
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("data sort value", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "data-sort-value.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := []Table{
			{
				Data: [][]string{
					{"Abu Dhabi, United Arab Emirates", "N/A"},
				},
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/data-sort-value", nil)

		svr.ServeHTTP(w, r)

		var got []Table
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusOK {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusOK)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}

func TestServerError(t *testing.T) {
	t.Run("row span error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "rowspanError.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := ServerError{
			Message: `strconv.Atoi: parsing "x": invalid syntax`,
			Metadata: map[string]interface{}{
				"CellNumber":         float64(0),
				"ResponseStatusCode": float64(500),
				"ResponseStatusText": "Internal Server Error",
				"RowNumber":          float64(1),
				"TableIndex":         float64(0),
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/rowspanError", nil)

		svr.ServeHTTP(w, r)

		var got ServerError
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusInternalServerError)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("row span error first table", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "rowspanError.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := ServerError{
			Message: `strconv.Atoi: parsing "x": invalid syntax`,
			Metadata: map[string]interface{}{
				"CellNumber":         float64(0),
				"ResponseStatusCode": float64(500),
				"ResponseStatusText": "Internal Server Error",
				"RowNumber":          float64(1),
				"TableIndex":         float64(0),
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/rowspanError", nil)
		params := r.URL.Query()
		params.Add("table", "0")
		r.URL.RawQuery = params.Encode()

		svr.ServeHTTP(w, r)

		var got ServerError
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusInternalServerError)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("col span error", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "colspanError.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := ServerError{
			Message: `strconv.Atoi: parsing "x": invalid syntax`,
			Metadata: map[string]interface{}{
				"CellNumber":         float64(1),
				"ResponseStatusCode": float64(500),
				"ResponseStatusText": "Internal Server Error",
				"RowNumber":          float64(1),
				"TableIndex":         float64(0),
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/colspanError", nil)

		svr.ServeHTTP(w, r)

		var got ServerError
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusInternalServerError {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusInternalServerError)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("wiki api not ok", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.WriteHeader(http.StatusRequestEntityTooLarge)
				w.Write([]byte("request entity too large"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := ServerError{
			Message: `request entity too large`,
			Metadata: map[string]interface{}{
				"Page":               "apiError",
				"ResponseStatusCode": float64(http.StatusRequestEntityTooLarge),
				"ResponseStatusText": http.StatusText(http.StatusRequestEntityTooLarge),
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/apiError", nil)

		svr.ServeHTTP(w, r)

		var got ServerError
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusRequestEntityTooLarge {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusRequestEntityTooLarge)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("bad table index param", func(t *testing.T) {
		ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			switch r.URL.Path {
			case "/":
				w.Write(getRespBody(t, "table.html"))
			default:
				t.Fatalf("path %s not supported", r.URL.Path)
			}
		}))
		defer ts.Close()

		want := ServerError{
			Message: `strconv.Atoi: parsing "x": invalid syntax`,
			Metadata: map[string]interface{}{
				"ResponseStatusCode": float64(http.StatusBadRequest),
				"ResponseStatusText": http.StatusText(http.StatusBadRequest),
			},
		}

		svr := NewServer(ts.URL)
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodGet, "/api/apiError", nil)
		params := r.URL.Query()
		params.Add("table", "x")
		r.URL.RawQuery = params.Encode()

		svr.ServeHTTP(w, r)

		var got ServerError
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusBadRequest {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusBadRequest)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})

	t.Run("non get method", func(t *testing.T) {
		want := ServerError{
			Message: "method POST not allowed",
			Metadata: map[string]interface{}{
				"ResponseStatusCode": float64(http.StatusBadRequest),
				"ResponseStatusText": http.StatusText(http.StatusBadRequest),
			},
		}

		svr := NewServer("")
		w := httptest.NewRecorder()
		r := httptest.NewRequest(http.MethodPost, "/api/post", nil)

		svr.ServeHTTP(w, r)

		var got ServerError
		err := json.Unmarshal(w.Body.Bytes(), &got)
		if err != nil {
			t.Fatal(err)
		}

		if w.Code != http.StatusBadRequest {
			t.Errorf("got code %d, expected %d", w.Code, http.StatusBadRequest)
		}

		if !reflect.DeepEqual(want, got) {
			t.Errorf("got %+v, want %+v", got, want)
		}
	})
}

func getRespBody(t *testing.T, file string) []byte {
	tables, err := os.ReadFile(fmt.Sprintf("%s/%s", "testdata", file))
	if err != nil {
		t.Fatal(err)
	}

	return tables
}

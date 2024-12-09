package server

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
	"time"

	"github.com/atye/wikitable2json/pkg/client"
	"github.com/atye/wikitable2json/pkg/client/status"
)

func TestServeHTTP_CacheMissGetMatrix(t *testing.T) {
	wantData := [][][]string{
		{
			{"test", "test"},
			{"test", "test"},
		},
	}

	tg := &mockTableGetter{getMatrix: wantData}
	sut := NewServer(tg, NewCache(10, 10*time.Second))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/page", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want code %d, got %d", http.StatusOK, w.Code)
	}

	if !tg.getMatrixCalled {
		t.Errorf("expected GetMatrix call")
	}

	var got [][][]string
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantData, got) {
		t.Errorf("expected %v, got %v", wantData, got)
	}
}

func TestServeHTTP_CacheMissGetMatrixVerbose(t *testing.T) {
	wantData := [][][]client.Verbose{
		{
			[]client.Verbose{
				{
					Text: "test0\ntest1",
					Links: []client.Link{
						{
							Href: "./test1",
							Text: "test1",
						},
					},
				},
				{
					Text: "test2",
				},
			},
		},
	}

	tg := &mockTableGetter{getMatrixVerbose: wantData}
	sut := NewServer(tg, NewCache(10, 10*time.Second))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/page?verbose=true", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want code %d, got %d", http.StatusOK, w.Code)
	}

	if !tg.getMatrixVerboseCalled {
		t.Errorf("expected GetMatrixVerbose call")
	}

	var got [][][]client.Verbose
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantData, got) {
		t.Errorf("expected %v, got %v", wantData, got)
	}
}

func TestServeHTTP_CacheMissGetKeyValue(t *testing.T) {
	wantData := [][]map[string]string{
		{
			{
				"Rank":    "1",
				"Account": "Alpha",
			},
		},
	}

	tg := &mockTableGetter{getKeyValue: wantData}
	sut := NewServer(tg, NewCache(10, 10*time.Second))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/page?keyRows=1", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want code %d, got %d", http.StatusOK, w.Code)
	}

	if !tg.getKeyValueCalled {
		t.Errorf("expected GetKeyValue call")
	}

	var got [][]map[string]string
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantData, got) {
		t.Errorf("expected %v, got %v", wantData, got)
	}
}

func TestServeHTTP_CacheMissGetKeyValueVerbose(t *testing.T) {
	wantData := [][]map[string]client.Verbose{
		{
			{
				"header1": {
					Text: "test",
				},
				"header2": {
					Text: "Bolivia, Plurinational State of",
					Links: []client.Link{
						{
							Text: "Bolivia, Plurinational State of",
							Href: "./Bolivia",
						},
					},
				},
			},
		},
	}

	tg := &mockTableGetter{getKeyValueVerbose: wantData}
	sut := NewServer(tg, NewCache(10, 10*time.Second))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/page?keyRows=1&verbose=true&cleanRef=true&brNewLine=true", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want code %d, got %d", http.StatusOK, w.Code)
	}

	if !tg.getKeyValueVerboseCalled {
		t.Errorf("expected GetKeyValue call")
	}

	var got [][]map[string]client.Verbose
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantData, got) {
		t.Errorf("expected %v, got %v", wantData, got)
	}
}

func TestServeHTTP_CacheHit(t *testing.T) {
	wantData := [][][]string{
		{
			{"test", "test"},
			{"test", "test"},
		},
	}

	tg := &mockTableGetter{getMatrix: wantData}
	sut := NewServer(tg, NewCache(10, 10*time.Second))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/page", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	tg.getMatrixCalled = false

	w = httptest.NewRecorder()
	r = httptest.NewRequest("GET", "/api/page", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusOK {
		t.Errorf("want code %d, got %d", http.StatusOK, w.Code)
	}

	if tg.getMatrixCalled {
		t.Errorf("expected GetMatrix not to be called")
	}

	var got [][][]string
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantData, got) {
		t.Errorf("expected %v, got %v", wantData, got)
	}
}

func TestServeHTTP_EmptyPage(t *testing.T) {
	wantData := status.Status{
		Message: "page value must be supplied in /api/{page}",
		Code:    http.StatusBadRequest,
	}

	sut := NewServer(&mockTableGetter{}, NewCache(10, 10*time.Second))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api", nil)
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want code %d, got %d", http.StatusBadRequest, w.Code)
	}

	var got status.Status
	err := json.Unmarshal(w.Body.Bytes(), &got)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(wantData, got) {
		t.Errorf("expected %v, got %v", wantData, got)
	}
}

func TestServeHTTP_InvalidQuery(t *testing.T) {
	sut := NewServer(&mockTableGetter{}, NewCache(10, 10*time.Second))

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/page?table=x", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusBadRequest {
		t.Errorf("want code %d, got %d", http.StatusBadRequest, w.Code)
	}
}

func TestServeHTTP_ClientError(t *testing.T) {
	err := status.NewStatus("error", http.StatusInternalServerError)

	cache := NewCache(10, 10*time.Second)
	tg := &mockTableGetter{getMatrix: nil, err: err}
	sut := NewServer(tg, cache)

	w := httptest.NewRecorder()
	r := httptest.NewRequest("GET", "/api/page", nil)
	r.SetPathValue("page", "page")
	sut.ServeHTTP(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Errorf("want code %d, got %d", http.StatusInternalServerError, w.Code)
	}

	if _, ok := cache.Get("page-en-all-false-0-false-false"); ok {
		t.Errorf("expected cache miss, got page-en-x-false-0-false-false")
	}
}

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
		{
			"en-all-true-2-true-true",
			queryValues{
				lang:      "en",
				cleanRef:  true,
				keyRows:   0,
				verbose:   true,
				brNewLine: true,
			},
			"test",
			"test-en-all-true-0-true-true",
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

type mockTableGetter struct {
	getMatrix                [][][]string
	getMatrixCalled          bool
	getMatrixVerbose         [][][]client.Verbose
	getMatrixVerboseCalled   bool
	getKeyValue              [][]map[string]string
	getKeyValueCalled        bool
	getKeyValueVerbose       [][]map[string]client.Verbose
	getKeyValueVerboseCalled bool
	userAgent                string
	err                      error
}

func (m *mockTableGetter) GetMatrix(ctx context.Context, page string, lang string, options ...client.TableOption) ([][][]string, error) {
	m.getMatrixCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return m.getMatrix, nil
}

func (m *mockTableGetter) GetMatrixVerbose(ctx context.Context, page string, lang string, options ...client.TableOption) ([][][]client.Verbose, error) {
	m.getMatrixVerboseCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return m.getMatrixVerbose, nil
}

func (m *mockTableGetter) GetKeyValue(ctx context.Context, page string, lang string, keyRows int, options ...client.TableOption) ([][]map[string]string, error) {
	m.getKeyValueCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return m.getKeyValue, nil
}

func (m *mockTableGetter) GetKeyValueVerbose(ctx context.Context, page string, lang string, keyRows int, options ...client.TableOption) ([][]map[string]client.Verbose, error) {
	m.getKeyValueVerboseCalled = true
	if m.err != nil {
		return nil, m.err
	}
	return m.getKeyValueVerbose, nil
}

func (m mockTableGetter) SetUserAgent(userAgent string) {
	m.userAgent = userAgent
}

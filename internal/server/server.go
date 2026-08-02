package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/atye/wikitable2json/pkg/client"
	"github.com/atye/wikitable2json/pkg/client/status"
)

const (
	defaultLang = "en"
)

type TableGetter interface {
	GetMatrix(ctx context.Context, page string, lang string, options ...client.TableOption) ([][][]string, error)
	GetMatrixVerbose(ctx context.Context, page string, lang string, options ...client.TableOption) ([][][]client.Verbose, error)
	GetKeyValue(ctx context.Context, page string, lang string, keyRows int, options ...client.TableOption) ([][]map[string]string, error)
	GetKeyValueVerbose(ctx context.Context, page string, lang string, keyRows int, options ...client.TableOption) ([][]map[string]client.Verbose, error)
}

type Server struct {
	client TableGetter
	cache  *Cache
}

func NewServer(client TableGetter, cache *Cache) (*Server, error) {
	if client == nil {
		return nil, fmt.Errorf("client not set")
	}

	if cache == nil {
		return nil, fmt.Errorf("cache not set")
	}
	return &Server{
		client: client,
		cache:  cache,
	}, nil
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page, ok := ctx.Value(pageKey).(string)
	if !ok {
		writeError(w, status.NewStatus("something went wrong. no page in request context", http.StatusInternalServerError))
		return
	}

	qv, ok := ctx.Value(queryKey).(queryValues)
	if !ok {
		writeError(w, status.NewStatus("something went wrong. no query values in request context", http.StatusInternalServerError))
		return
	}

	key, err := buildCacheKey(page, qv)
	if err != nil {
		writeError(w, status.NewStatus(err.Error(), http.StatusInternalServerError))
		return
	}

	data, ok := s.cache.Get(key)
	if ok {
		err = json.NewEncoder(w).Encode(data)
		if err != nil {
			writeError(w, status.NewStatus(err.Error(), http.StatusInternalServerError))
			return
		}
		return
	}

	opts := []client.TableOption{
		client.WithTables(qv.tables...),
		client.WithSections(qv.sections...),
	}
	if qv.cleanRef {
		opts = append(opts, client.WithCleanReferences())
	}
	if qv.brNewLine {
		opts = append(opts, client.WithBRNewLine())
	}

	var resp interface{}
	if qv.keyRows >= 1 {
		if qv.verbose {
			resp, err = s.client.GetKeyValueVerbose(ctx, page, qv.lang, qv.keyRows, opts...)
		} else {
			resp, err = s.client.GetKeyValue(ctx, page, qv.lang, qv.keyRows, opts...)
		}
	} else {
		if qv.verbose {
			resp, err = s.client.GetMatrixVerbose(ctx, page, qv.lang, opts...)
		} else {
			resp, err = s.client.GetMatrix(ctx, page, qv.lang, opts...)
		}
	}
	if err != nil {
		writeError(w, err)
		return
	}

	defer func() {
		_ = s.cache.Add(key, resp)
	}()

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(w, status.NewStatus(err.Error(), http.StatusInternalServerError))
		return
	}
}

type queryValues struct {
	lang      string
	tables    []int
	sections  []string
	cleanRef  bool
	keyRows   int
	verbose   bool
	brNewLine bool
}

func parseParameters(r *http.Request) (queryValues, error) {
	var qv queryValues
	qv.lang = defaultLang

	params := r.URL.Query()

	if v := params.Get("lang"); v != "" {
		qv.lang = v
	}

	if v, ok := params["table"]; ok {
		for _, table := range v {
			t, err := strconv.Atoi(table)
			if err != nil {
				return queryValues{}, status.NewStatus(err.Error(), http.StatusBadRequest)
			}
			qv.tables = append(qv.tables, t)
		}
	}

	if v, ok := params["section"]; ok {
		qv.sections = v
	}

	if v := params.Get("cleanRef"); v == "true" {
		qv.cleanRef = true
	}

	if v := params.Get("verbose"); v == "true" {
		qv.verbose = true
	}

	if v := params.Get("brNewLine"); v == "true" {
		qv.brNewLine = true
	}

	if v := params.Get("keyRows"); v != "" {
		n, err := strconv.Atoi(v)
		if err != nil {
			return queryValues{}, status.NewStatus(err.Error(), http.StatusBadRequest)
		}

		if n < 1 {
			return queryValues{}, status.NewStatus("keyRows must be at least 1", http.StatusBadRequest)
		}

		qv.keyRows = n
	}

	return qv, nil
}

func buildCacheKey(page string, qv queryValues) (string, error) {
	key := cacheKey{
		Page:      page,
		Lang:      qv.lang,
		Tables:    qv.tables,
		Sections:  qv.sections,
		CleanRef:  qv.cleanRef,
		KeyRows:   qv.keyRows,
		Verbose:   qv.verbose,
		BrNewLine: qv.brNewLine,
	}

	b, err := json.Marshal(key)
	if err != nil {
		return "", fmt.Errorf("marshalling cache key: %v", err)
	}

	return string(b), nil
}

func writeError(w http.ResponseWriter, err error) {
	var s status.Status
	if errors.As(err, &s) {
		code := s.Code
		if code == 0 {
			code = http.StatusInternalServerError
		}
		w.WriteHeader(code)

		err = json.NewEncoder(w).Encode(s)
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshaling error response: %v", err), http.StatusInternalServerError)
		}
	} else {
		w.WriteHeader(http.StatusInternalServerError)
		err = json.NewEncoder(w).Encode(status.NewStatus(err.Error(), http.StatusInternalServerError))
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshaling error response: %v", err), http.StatusInternalServerError)
		}
	}
}

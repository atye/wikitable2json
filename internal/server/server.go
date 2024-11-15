package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/atye/wikitable2json/internal/status"
	"github.com/atye/wikitable2json/pkg/client"
)

const (
	defaultLang = "en"
)

type Server struct {
	client client.TableGetter
}

func NewServer(client client.TableGetter) *Server {
	return &Server{
		client: client,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	page := r.PathValue("page")

	qv, err := parseParameters(r)
	if err != nil {
		writeError(w, err)
		return
	}

	opts := []client.TableOption{
		client.WithTables(qv.tables...),
	}
	if qv.cleanRef {
		opts = append(opts, client.WithCleanReferences())
	}
	if qv.brNewLine {
		opts = append(opts, client.WithBRNewLine())
	}

	s.client.SetUserAgent(r.Header.Get("User-Agent"))

	var resp interface{}
	if qv.keyRows >= 1 {
		if qv.verbose {
			resp, err = s.client.GetKeyValueVerbose(ctx, page, qv.lang, qv.keyRows, opts...)
			if err != nil {
				writeError(w, err)
				return
			}
		} else {
			resp, err = s.client.GetKeyValue(ctx, page, qv.lang, qv.keyRows, opts...)
			if err != nil {
				writeError(w, err)
				return
			}
		}
	} else {
		if qv.verbose {
			resp, err = s.client.GetMatrixVerbose(ctx, page, qv.lang, opts...)
			if err != nil {
				writeError(w, err)
				return
			}
		} else {
			resp, err = s.client.GetMatrix(ctx, page, qv.lang, opts...)
			if err != nil {
				writeError(w, err)
				return
			}
		}
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(w, status.NewStatus(err.Error(), http.StatusInternalServerError))
		return
	}
}

type queryValues struct {
	lang      string
	tables    []int
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

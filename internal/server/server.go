package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

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
	if r.Method != http.MethodGet {
		writeError(w, status.NewStatus(fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed))
		return
	}

	ctx := r.Context()

	page := strings.TrimPrefix(r.URL.Path, "/api/")

	qv, err := parseParameters(r)
	if err != nil {
		writeError(w, err)
		return
	}

	s.client.SetUserAgent(r.Header.Get("User-Agent"))

	var resp interface{}
	if qv.keyRows >= 1 {
		resp, err = s.client.GetTablesKeyValue(ctx, page, qv.lang, qv.cleanRef, qv.keyRows, qv.tables...)
		if err != nil {
			writeError(w, err)
			return
		}
	} else {
		resp, err = s.client.GetTablesMatrix(ctx, page, qv.lang, qv.cleanRef, qv.tables...)
		if err != nil {
			writeError(w, err)
			return
		}
	}

	err = json.NewEncoder(w).Encode(resp)
	if err != nil {
		writeError(w, status.NewStatus(err.Error(), http.StatusInternalServerError))
		return
	}
}

type queryValues struct {
	lang     string
	tables   []int
	cleanRef bool
	keyRows  int
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
		if s.Code != 0 {
			w.WriteHeader(s.Code)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
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

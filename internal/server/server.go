package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/atye/wikitable-api/internal/status"
)

const (
	defaultLang   = "en"
	defaultFormat = "matrix"
)

type WikiAPI interface {
	GetPageData(ctx context.Context, page, lang string) (io.ReadCloser, error)
}

type Server struct {
	wiki WikiAPI
}

func NewServer(wiki WikiAPI) *Server {
	return &Server{
		wiki: wiki,
	}
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeError(w, status.NewStatus(fmt.Sprintf("method %s not allowed", r.Method), http.StatusMethodNotAllowed))
		return
	}

	page := strings.TrimPrefix(r.URL.Path, "/api/")
	lang, format, tables, cleanRef, err := parseParameters(r)
	if err != nil {
		writeError(w, err)
		return
	}

	reader, err := s.wiki.GetPageData(r.Context(), page, lang)
	if err != nil {
		writeError(w, err)
		return
	}
	defer reader.Close()

	input := parseOptions{
		tables:   tables,
		format:   format,
		cleanRef: cleanRef,
	}

	resp, err := parse(r.Context(), reader, input)
	if err != nil {
		writeError(w, err)
		return
	}

	b, err := json.Marshal(resp)
	if err != nil {
		writeError(w, status.NewStatus(err.Error(), http.StatusInternalServerError))
		return
	}

	fmt.Fprintf(w, "%s", b)
}

func parseParameters(r *http.Request) (lang string, format string, tables []int, cleanRef bool, e error) {
	params := r.URL.Query()
	lang = defaultLang
	if v := params.Get("lang"); v != "" {
		lang = v
	}

	if v, ok := params["table"]; ok {
		for _, table := range v {
			t, err := strconv.Atoi(table)
			if err != nil {
				e = status.NewStatus(err.Error(), http.StatusBadRequest)
				return
			}
			tables = append(tables, t)
		}
	}

	format = defaultFormat
	if v := params.Get("format"); v != "" {
		format = v
	}

	if v := params.Get("cleanRef"); v == "true" {
		cleanRef = true
	}

	return
}

func writeError(w http.ResponseWriter, err error) {
	var s status.Status
	if errors.As(err, &s) {
		b, err := json.Marshal(s)
		if err != nil {
			http.Error(w, fmt.Sprintf("error marshaling error response: %v", err), http.StatusInternalServerError)
			return
		}

		if s.Code != 0 {
			w.WriteHeader(s.Code)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}

		fmt.Fprintf(w, "%s", b)
		return
	}

	b, err := json.Marshal(status.NewStatus(err.Error(), http.StatusInternalServerError))
	if err != nil {
		http.Error(w, fmt.Sprintf("error marshaling error response: %v", err), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%s", b)
}

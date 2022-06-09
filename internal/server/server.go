package server

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable2json/internal/status"
)

const (
	defaultLang   = "en"
	defaultFormat = "matrix"
)

var (
	classes = []string{
		"table.wikitable",
		"table.standard",
		"table.toccolours",
	}
)

type WikiAPI interface {
	GetPageBytes(ctx context.Context, page, lang, userAgent string) ([]byte, error)
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

	qv, err := parseParameters(r)
	if err != nil {
		writeError(w, err)
		return
	}

	data, err := s.wiki.GetPageBytes(r.Context(), page, qv.lang, r.Header.Get("User-Agent"))
	if err != nil {
		writeError(w, err)
		return
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(data))
	if err != nil {
		writeError(w, err)
		return
	}

	doc.Find(".mw-empty-elt").Remove()
	tables := doc.Find(strings.Join(classes, ", "))

	if qv.cleanRef {
		cleanReferences(tables)
	}

	opts := parseOptions{
		tables:  qv.tables,
		keyrows: qv.keyRows,
	}

	resp, err := parse(r.Context(), tables, opts)
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

func cleanReferences(tables *goquery.Selection) {
	tables.Find(".reference").Remove()

	tables.Find("sup").Each(func(_ int, s *goquery.Selection) {
		s.Find("a").Each(func(_ int, anchor *goquery.Selection) {
			if v, ok := anchor.Attr("title"); ok {
				if v == "Wikipedia:Citation needed" {
					s.Remove()
				}
			}
		})
	})
}

type queryValues struct {
	lang     string
	tables   []int
	cleanRef bool
	keyRows  int
}

func parseParameters(r *http.Request) (queryValues, error) {
	var qv queryValues

	params := r.URL.Query()
	qv.lang = defaultLang
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

package v2

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	defaultLang = "en"
)

var (
	baseURL = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
)

type Server struct{}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		writeServerError(w, generalErr(fmt.Errorf("method %s not allowed", r.Method), http.StatusBadRequest))
		return
	}

	params := r.URL.Query()

	lang := defaultLang
	if v := params.Get("lang"); v != "" {
		lang = v
	}

	doc, err := s.getDocument(r.Context(), strings.TrimPrefix(r.URL.Path, "/api/"), lang)
	if err != nil {
		writeServerError(w, err)
		return
	}

	var tableParams []string
	if v, ok := params["table"]; ok {
		tableParams = v
	}

	tables, err := parseTables(r.Context(), doc.Find("table.wikitable"), tableParams)
	if err != nil {
		writeServerError(w, err)
		return
	}

	resp, err := json.Marshal(tables)
	if err != nil {
		writeServerError(w, generalErr(err, http.StatusInternalServerError))
		return
	}

	fmt.Fprintf(w, "%s", resp)
}

func (s *Server) getDocument(ctx context.Context, page, lang string) (*goquery.Document, error) {
	resp, err := s.getWikiAPIResponse(ctx, page, lang)
	if err != nil {
		return nil, err
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	doc.Find(".mw-empty-elt").Remove()
	return doc, nil
}

type wikiApiError struct {
	statusCode int
	message    string
	page       string
}

func (e *wikiApiError) Error() string {
	return e.message
}

func (s *Server) getWikiAPIResponse(ctx context.Context, page, lang string) (*http.Response, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(baseURL, lang, url.QueryEscape(page)), nil)
	if err != nil {
		return nil, generalErr(err, http.StatusInternalServerError)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, generalErr(err, http.StatusInternalServerError)
	}
	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, wikiAPIErr(&wikiApiError{statusCode: resp.StatusCode, page: page, message: fmt.Sprintf("failed to read wikipedia API response body: %v", err.Error())})
		}
		return nil, wikiAPIErr(&wikiApiError{statusCode: resp.StatusCode, page: page, message: string(body)})
	}
	return resp, nil
}

func writeServerError(w http.ResponseWriter, err error) {
	var svrErr *ServerError
	if errors.As(err, &svrErr) {
		if codeValue, ok := svrErr.Metadata["ResponseStatusCode"]; ok {
			if code, ok := codeValue.(int); ok {
				w.WriteHeader(code)
			} else {
				w.WriteHeader(http.StatusInternalServerError)
			}
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		log.Printf("%+v", svrErr)
		fmt.Fprintf(w, "%v", svrErr)
		return
	}
	genErr := generalErr(err, http.StatusInternalServerError)
	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintf(w, "%v", genErr)
}

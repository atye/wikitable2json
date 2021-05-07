package v2

import (
	"context"
	"fmt"
	"io/ioutil"
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
		//
	}

	params := r.URL.Query()

	lang := defaultLang
	if v := params.Get("lang"); v != "" {
		lang = v
	}

	_, _ = s.getDocument(r.Context(), strings.TrimPrefix(r.URL.Path, "/api/"), lang)

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
		return nil, generalErr(err)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, generalErr(err)
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, &wikiApiError{statusCode: resp.StatusCode, page: page, message: fmt.Sprintf("failed to read wikipedia API response body: %v", err.Error())}
		}
		return nil, &wikiApiError{statusCode: resp.StatusCode, page: page, message: string(body)}
	}
	return resp, nil
}

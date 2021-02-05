package service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable-api/internal/service/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Service struct {
	HTTPGet func(string) (*http.Response, error)
	pb.UnimplementedWikipediaTableJSONAPIServer
}

const (
	defaultLang = "en"
)

var (
	baseURL = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
)

func (s *Service) GetTables(ctx context.Context, req *pb.TablesRequest) (*pb.TablesResponse, error) {
	doc, err := s.getDocument(ctx, req)
	if err != nil {
		var apiErr *wikiApiError
		if errors.As(err, &apiErr) {
			return nil, wikiAPIStatusErr(apiErr)
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to retrieve wikipedia API response document: %v", err))
	}
	resp, err := parseTables(ctx, doc.Find("table.wikitable"), req.Table)
	if err != nil {
		var ptErr *parseTableError
		if errors.As(err, &ptErr) {
			return nil, tableParseStatusErr(ptErr)
		}
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to parse a table: %v", err))
	}
	return resp, nil
}

func (s *Service) getDocument(ctx context.Context, req *pb.TablesRequest) (*goquery.Document, error) {
	lang := defaultLang
	if req.Lang != "" {
		lang = req.Lang
	}
	resp, err := s.getWikiAPIResponse(ctx, req.Page, lang)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
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
	return ""
}

func (s *Service) getWikiAPIResponse(ctx context.Context, page, lang string) (*http.Response, error) {
	resp, err := s.HTTPGet(fmt.Sprintf(baseURL, lang, url.QueryEscape(page)))
	if err != nil {
		return nil, status.Error(codes.Internal, fmt.Sprintf("failed to make http request to the wikipedia API: %v", err))
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, &wikiApiError{statusCode: resp.StatusCode, page: page, message: fmt.Sprintf("failed to read wikipedia API response body: %v", err)}
		}
		return nil, &wikiApiError{statusCode: resp.StatusCode, page: page, message: string(body)}
	}
	return resp, nil
}

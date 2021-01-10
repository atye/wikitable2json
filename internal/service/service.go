package service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable-api/internal/service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Service struct {
	HTTPGet func(string) (*http.Response, error)
	pb.UnimplementedWikipediaTableJSONAPIServer
}

type wikiApiError struct {
	statusCode int
	message    string
	page       string
}

func (e *wikiApiError) Error() string {
	return ""
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
			return nil, wikiAPIRespNotOKStatusErr(apiErr)
		}
		return nil, getGeneralStatusErr(err, "something went wrong retrieving the wikipedia API response document")
	}
	resp, err := parseTables(ctx, doc.Find("table.wikitable"), req.Table)
	if err != nil {
		var ptErr *parseTableError
		if errors.As(err, &ptErr) {
			return nil, tableParseStatusErr(ptErr)
		}
		return nil, getGeneralStatusErr(err, "something went wrong parsing a table")
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
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}
	doc.Find(".mw-empty-elt").Remove()
	return doc, nil
}

func (s *Service) getWikiAPIResponse(ctx context.Context, page, lang string) (*http.Response, error) {
	resp, err := s.HTTPGet(fmt.Sprintf(baseURL, lang, url.QueryEscape(page)))
	if err != nil {
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusInternalServerError)))
		return nil, err
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

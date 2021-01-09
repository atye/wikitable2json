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
		return nil, getGeneralStatusErr(err)
	}
	resp, err := parseTables(ctx, doc.Find("table.wikitable"), int32ToInt(req.Table))
	if err != nil {
		var ptErr *parseTableError
		if errors.As(err, &ptErr) {
			return nil, tableParseStatusErr(ptErr)
		}
		panic("unsupported parseTable error type")
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

func int32ToInt(input []int32) []int {
	if input != nil && len(input) > 0 {
		output := make([]int, len(input))
		for i, value := range input {
			output[i] = int(value)
		}
		return output
	}
	return []int{}
}

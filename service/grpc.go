package service

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"strconv"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable-api/service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Service struct {
	HTTPGet func(string) (*http.Response, error)
}

const (
	baseURL     = "wikipedia.org/api/rest_v1/page/html"
	defaultLang = "en"
)

var (
	ErrWikipediaRestAPINotOk = errors.New("Wikiedia API response not OK")
)

func (s *Service) GetTables(ctx context.Context, req *pb.GetTablesRequest) (*pb.GetTablesResponse, error) {
	doc, statusCode, err := getDocument(req, s.HTTPGet)
	if err != nil {
		if errors.Is(err, ErrWikipediaRestAPINotOk) {
			headerErr := grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(statusCode)))
			if headerErr != nil {
				return nil, headerErr
			}
			return nil, fmt.Errorf("wikipedia API response not 200/OK - got status code: %d", statusCode)
		}
		return nil, err
	}

	n := make([]int, len(req.N))
	for i, reqN := range req.N {
		index, err := strconv.Atoi(reqN)
		if err != nil {
			headerErr := grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusBadRequest)))
			if headerErr != nil {
				return nil, headerErr
			}
			return nil, fmt.Errorf("table index (n) should be a number - got %s", reqN)
		}
		n[i] = index
	}

	resp, err := parseTables(ctx, doc.Find("table.wikitable"), n)
	if err != nil {
		headerErr := grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusInternalServerError)))
		if headerErr != nil {
			return nil, headerErr
		}
		return nil, err
	}

	return resp, nil
}

func getDocument(req *pb.GetTablesRequest, get func(string) (*http.Response, error)) (*goquery.Document, int, error) {
	lang := defaultLang
	if req.Lang != "" {
		lang = req.Lang
	}

	resp, err := get(fmt.Sprintf("https://%s.%s/%s", lang, baseURL, url.QueryEscape(req.Page)))
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, resp.StatusCode, ErrWikipediaRestAPINotOk
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, err
	}

	// remove empty/hidden elements
	doc.Find(".mw-empty-elt").Remove()

	return doc, resp.StatusCode, err
}

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
	"github.com/atye/wikitable-api/service/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type Service struct {
	HttpGet func(string) (*http.Response, error)
	pb.UnimplementedWikiTableJSONAPIServer
}

const (
	defaultLang = "en"
)

var (
	baseURL = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
)

func (s *Service) GetTables(ctx context.Context, req *pb.TablesRequest) (*pb.TablesResponse, error) {
	doc, err := getDocument(ctx, req, s.HttpGet)
	if err != nil {
		return nil, err
	}
	resp, err := parseTables(ctx, doc.Find("table.wikitable"), int32ToInt(req.Table))
	if err != nil {
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusInternalServerError)))
		return nil, err
	}
	return resp, nil
}

func getDocument(ctx context.Context, req *pb.TablesRequest, httpGet func(string) (*http.Response, error)) (*goquery.Document, error) {
	lang := defaultLang
	if req.Lang != "" {
		lang = req.Lang
	}
	resp, err := httpGet(fmt.Sprintf(baseURL, lang, url.QueryEscape(req.Page)))
	if err != nil {
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusServiceUnavailable)))
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusInternalServerError)))
			return nil, err
		}
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(resp.StatusCode)))
		return nil, errors.New(string(body))
	}
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusInternalServerError)))
		return nil, err
	}
	doc.Find(".mw-empty-elt").Remove()
	return doc, nil
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

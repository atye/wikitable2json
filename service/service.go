package service

import (
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable-api/service/pb"
	"golang.org/x/sync/errgroup"
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
	doc, err := s.getDocument(ctx, req, s.HttpGet)
	if err != nil {
		var apiErr *wikiApiError
		if errors.As(err, &apiErr) {
			grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(apiErr.statusCode)))
			return nil, wikiAPIStatusErr(apiErr)
		}
		return nil, getDocumentStatusErr(err)
	}
	resp, err := parseTables(ctx, doc.Find("table.wikitable"), int32ToInt(req.Table))
	if err != nil {
		return nil, tableParseStatusErr(err)
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

func (s *Service) getDocument(ctx context.Context, req *pb.TablesRequest, httpGet func(string) (*http.Response, error)) (*goquery.Document, error) {
	lang := defaultLang
	if req.Lang != "" {
		lang = req.Lang
	}
	resp, err := s.getWikiAPIResponse(ctx, fmt.Sprintf(baseURL, lang, url.QueryEscape(req.Page)))
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

func (s *Service) getWikiAPIResponse(ctx context.Context, url string) (*http.Response, error) {
	resp, err := s.HttpGet(url)
	if err != nil {
		grpc.SetHeader(ctx, metadata.Pairs("x-http-code", strconv.Itoa(http.StatusInternalServerError)))
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, &wikiApiError{statusCode: resp.StatusCode, message: fmt.Sprintf("failed to read wikipedia API response body: %v", err.Error())}
		}
		return nil, &wikiApiError{statusCode: resp.StatusCode, message: string(body)}
	}
	return resp, nil
}

func parseTables(ctx context.Context, wikiTableSelection *goquery.Selection, tableIndices []int) (*pb.TablesResponse, error) {
	var eg errgroup.Group
	switch len(tableIndices) {
	case 0:
		resp := &pb.TablesResponse{
			Tables: make([]*pb.Table, len(wikiTableSelection.Nodes)),
		}
		wikiTableSelection.Each(func(i int, selection *goquery.Selection) {
			eg.Go(func() error {
				table, err := parseTable(selection)
				if err != nil {
					return err
				}
				resp.Tables[i] = table
				return nil
			})
		})
		err := eg.Wait()
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		resp := &pb.TablesResponse{
			Tables: make([]*pb.Table, len(tableIndices)),
		}
		for i, tableIndex := range tableIndices {
			i := i
			tableIndex := tableIndex
			eg.Go(func() error {
				table, err := parseTable(wikiTableSelection.Eq(tableIndex))
				if err != nil {
					return err
				}
				resp.Tables[i] = table
				return nil
			})
		}
		err := eg.Wait()
		if err != nil {
			return nil, err
		}
		return resp, nil
	}
}

func parseTable(tableSelection *goquery.Selection) (*pb.Table, error) {
	table := &pb.Table{
		Rows: make(map[int64]*pb.Row),
	}
	table.Caption = tableSelection.Find("caption").Text()
	var err error
	// for each row in the table
	tableSelection.Find("tr").EachWithBreak(func(rowNum int, s *goquery.Selection) bool {
		// find all th and td elements in the row
		s.Find("th, td").EachWithBreak(func(cellNum int, s *goquery.Selection) bool {
			rowSpan := 1
			colSpan := 1
			// get the rowspan and colspan attributes
			if attr := s.AttrOr("rowspan", ""); attr != "" {
				rowSpan, err = strconv.Atoi(attr)
				if err != nil {
					return false
				}
			}
			if attr := s.AttrOr("colspan", ""); attr != "" {
				colSpan, err = strconv.Atoi(attr)
				if err != nil {
					return false
				}
			}
			// loop through the spans and populate table columns
			for i := 0; i < rowSpan; i++ {
				for j := 0; j < colSpan; j++ {
					row := rowNum + i
					if _, ok := table.Rows[int64(row)]; !ok {
						table.Rows[int64(row)] = &pb.Row{
							Columns: make(map[int64]string),
						}
					}
					nextAvailableCell := 0
					columns := table.Rows[int64(row)].Columns
					// check if column already has a value from a previous rowspan so we don't overrwite it
					// loop until we get an availalbe column
					// https://en.wikipedia.org/wiki/Help:Table#Combined_use_of_COLSPAN_and_ROWSPAN
					for columns[int64(cellNum+j+nextAvailableCell)] != "" {
						nextAvailableCell++
					}
					columns[int64(cellNum+j+nextAvailableCell)] = parseText(s.Text())
				}
			}
			return true
		})
		if err != nil {
			return false
		}
		return true
	})
	if err != nil {
		return nil, err
	}
	return table, nil
}

func parseText(s string) string {
	return strings.TrimSpace(s)
}

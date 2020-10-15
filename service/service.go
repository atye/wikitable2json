package service

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable-api/service/pb"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
)

const (
	baseURL     = "wikipedia.org/api/rest_v1/page/html"
	defaultLang = "en"
)

type Service struct{}

func (s *Service) GetTables(ctx context.Context, req *pb.GetTablesRequest) (*pb.GetTablesResponse, error) {
	doc, err := getDocument(req)
	if err != nil {
		return nil, err
	}

	wikiTableSelection := doc.Find("table.wikitable")
	var eg errgroup.Group

	switch len(req.N) {
	case 0:
		resp := &pb.GetTablesResponse{
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
		resp := &pb.GetTablesResponse{
			Tables: make([]*pb.Table, len(req.N)),
		}

		for i, n := range req.N {
			i := i
			n := n
			eg.Go(func() error {
				var index int
				index, err = strconv.Atoi(n)
				if err != nil {
					return err
				}

				table, err := parseTable(wikiTableSelection.Eq(index))
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

func getDocument(req *pb.GetTablesRequest) (*goquery.Document, error) {
	lang := defaultLang
	if req.Lang != "" {
		lang = req.Lang
	}

	url := fmt.Sprintf("https://%s.%s/%s", lang, baseURL, url.QueryEscape(req.Page))

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, grpc.Errorf(codes.Unknown, fmt.Sprintf("failed to get %s with status: %d", url, resp.StatusCode))
	}

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// remove empty/hidden elements
	doc.Find(".mw-empty-elt").Remove()

	return doc, err
}

func parseText(s string) string {
	s = strings.TrimSpace(s)
	return s
}

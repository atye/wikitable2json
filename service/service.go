package service

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable-api/service/pb"
)

const (
	baseURL = "https://en.wikipedia.org/api/rest_v1/page/html"
)

type Service struct {
}

func (s *Service) GetTables(ctx context.Context, req *pb.GetTablesRequest) (*pb.GetTablesResponse, error) {
	doc, err := getDocument(req.Page)
	if err != nil {
		return &pb.GetTablesResponse{}, err
	}

	errCh := make(chan error)
	resp := &pb.GetTablesResponse{}

	go func() {
		doc.Find("table.wikitable").Each(func(_ int, s *goquery.Selection) {
			table, err := parseTable(s)
			if err != nil {
				errCh <- err
			}

			resp.Tables = append(resp.Tables, table)
		})
		errCh <- nil
	}()

	err = <-errCh
	if err != nil {
		return &pb.GetTablesResponse{}, err
	}

	return resp, nil
}

func (s *Service) GetTable(ctx context.Context, req *pb.GetTableRequest) (*pb.Table, error) {
	doc, err := getDocument(req.Page)
	if err != nil {
		return &pb.Table{}, err
	}

	index, err := strconv.Atoi(req.N)
	if err != nil {
		return &pb.Table{}, err
	}

	table, err := parseTable(doc.Find("table.wikitable").Eq(index))
	if err != nil {
		return &pb.Table{}, err
	}

	return table, nil
}

func parseTable(tableSelection *goquery.Selection) (*pb.Table, error) {
	errCh := make(chan error)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// get the table rows, initialize a table, and get the table caption
	rows := tableSelection.Find("tr")

	table := initTable(rows)
	table.Caption = tableSelection.Find("caption").Text()

	go func() {
		// for each row in the table
		rows.Each(func(rowNum int, s *goquery.Selection) {
			// find all th and td elements in the row
			s.Find("th, td").Each(func(cellNum int, s *goquery.Selection) {
				// check if ctx was cancelled from a previous func getting row or col span attribute
				// if yes, just return to speed up this go routine cleanup
				if ctx.Err() != nil {
					return
				}

				var err error
				rowSpan := 1
				colSpan := 1

				// get the rowspan and colspan attributes
				// cancel the context if we get an error so we can quickly return from future .Each funcs
				if attr := s.AttrOr("rowspan", ""); attr != "" {
					rowSpan, err = strconv.Atoi(attr)
					if err != nil {
						errCh <- err
						cancel()
					}
				}

				if attr := s.AttrOr("colspan", ""); attr != "" {
					colSpan, err = strconv.Atoi(attr)
					if err != nil {
						errCh <- err
						cancel()
					}
				}

				// loop through the spans and populate table columns
				for i := 0; i < rowSpan; i++ {
					for j := 0; j < colSpan; j++ {
						nextAvailableCell := 0
						columns := table.Rows[int64(rowNum+i)].Columns

						// check if column already has a value from a previous rowspan so we don't overrwite it
						// loop until we get an availalbe column
						// https://en.wikipedia.org/wiki/Help:Table#Combined_use_of_COLSPAN_and_ROWSPAN
						for columns[int64(cellNum+j+nextAvailableCell)] != "" {
							nextAvailableCell++
						}
						columns[int64(cellNum+j+nextAvailableCell)] = parseText(s.Text())
					}
				}

			})
		})
		errCh <- nil
	}()

	err := <-errCh
	if err != nil {
		return nil, err
	}

	return table, nil
}

func getDocument(page string) (*goquery.Document, error) {
	resp, err := http.Get(fmt.Sprintf("%s/%s", baseURL, page))
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		return nil, err
	}

	// remove empty/hidden elements
	doc.Find(".mw-empty-elt").Remove()

	return doc, err
}

func initTable(rows *goquery.Selection) *pb.Table {
	table := &pb.Table{
		Rows: make(map[int64]*pb.Row),
	}

	for row := 0; row < rows.Length(); row++ {
		table.Rows[int64(row)] = &pb.Row{
			Columns: make(map[int64]string),
		}
	}

	return table
}

func parseText(s string) string {
	s = strings.TrimSpace(s)
	return s
}

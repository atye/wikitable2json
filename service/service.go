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
	baseURL = "wikipedia.org/api/rest_v1/page/html"
)

type Service struct{}

func (s *Service) GetTables(ctx context.Context, req *pb.GetTablesRequest) (*pb.GetTablesResponse, error) {
	var err error

	doc, err := getDocument(req)
	if err != nil {
		return &pb.GetTablesResponse{}, err
	}

	resp := &pb.GetTablesResponse{}
	var table *pb.Table

	switch len(req.N) {
	case 0:
		doc.Find("table.wikitable").EachWithBreak(func(_ int, s *goquery.Selection) bool {
			table, err = parseTable(s)
			if err != nil {
				return false
			}

			resp.Tables = append(resp.Tables, table)
			return true
		})
	default:
		for _, n := range req.N {
			var index int

			index, err = strconv.Atoi(n)
			if err != nil {
				return &pb.GetTablesResponse{}, err
			}

			table, err = parseTable(doc.Find("table.wikitable").Eq(index))
			if err != nil {
				return &pb.GetTablesResponse{}, err
			}

			resp.Tables = append(resp.Tables, table)
		}
	}

	if err != nil {
		return &pb.GetTablesResponse{}, err
	}

	return resp, nil
}

func parseTable(tableSelection *goquery.Selection) (*pb.Table, error) {
	// get the table rows, initialize a table, and get the table caption
	rows := tableSelection.Find("tr")

	table := initTable(rows)
	table.Caption = tableSelection.Find("caption").Text()

	var err error

	// for each row in the table
	rows.EachWithBreak(func(rowNum int, s *goquery.Selection) bool {
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

			return true
		})

		if err != nil {
			return false
		}

		return true
	})

	if err != nil {
		return &pb.Table{}, err
	}

	return table, nil
}

func getDocument(req *pb.GetTablesRequest) (*goquery.Document, error) {
	lang := "en"
	if req.Lang != "" {
		lang = req.Lang
	}

	resp, err := http.Get(fmt.Sprintf("https://%s.%s/%s", lang, baseURL, req.Page))
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

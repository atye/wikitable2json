package server

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable-api/internal/status"
	"golang.org/x/sync/errgroup"
)

var (
	classes = []string{
		"table.wikitable",
		"table.standard",
		"table.toccolours",
	}
)

func parse(ctx context.Context, r io.Reader, tables []int, format string) (interface{}, error) {
	tableSelection, err := getTableSelection(r)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	var eg errgroup.Group
	switch len(tables) {
	case 0:
		resp := make([]interface{}, tableSelection.Length())
		tableSelection.Each(func(i int, selection *goquery.Selection) {
			eg.Go(func() error {
				td, err := parseTable(selection, i)
				if err != nil {
					return err
				}

				f, err := toFormat(format, td, i)
				if err != nil {
					return err
				}

				resp[i] = f
				return nil
			})
		})
		err := eg.Wait()
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		resp := make([]interface{}, len(tables))
		for i, tableIndex := range tables {
			i := i
			tableIndex := tableIndex
			eg.Go(func() error {
				td, err := parseTable(tableSelection.Eq(tableIndex), tableIndex)
				if err != nil {
					return err
				}

				f, err := toFormat(format, td, i)
				if err != nil {
					return err
				}

				resp[i] = f
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

func getTableSelection(r io.Reader) (*goquery.Selection, error) {
	doc, err := goquery.NewDocumentFromReader(r)
	if err != nil {
		return nil, err
	}
	doc.Find(".mw-empty-elt").Remove()

	return doc.Find(strings.Join(classes, ", ")), nil
}

func parseTable(tableSelection *goquery.Selection, tableIndex int) (verbose, error) {
	td := make(verbose)

	errorStatus := status.Status{}
	var err error
	// for each row in the table
	tableSelection.Find("tr").EachWithBreak(func(rowNum int, s *goquery.Selection) bool {
		// find all th and td elements in the row
		s.Find("td, th").EachWithBreak(func(cellNum int, s *goquery.Selection) bool {
			rowSpan := 1
			colSpan := 1
			// get the rowspan and colspan attributes
			if attr := s.AttrOr("rowspan", ""); attr != "" {
				rowSpanTexts := strings.Split(attr, " ")
				if len(rowSpanTexts) > 0 {
					rowSpan, err = getSpan(rowSpanTexts)
					if err != nil {
						errorStatus = status.NewStatus(err.Error(), http.StatusInternalServerError, status.WithDetails(status.Details{
							status.TableIndex:   tableIndex,
							status.RowNumber:    rowNum,
							status.ColumnNumber: cellNum,
						}))
						return false
					}
				}
			}
			if attr := s.AttrOr("colspan", ""); attr != "" {
				colSpanTexts := strings.Split(attr, " ")
				if len(colSpanTexts) > 0 {
					colSpan, err = getSpan(colSpanTexts)
					if err != nil {
						errorStatus = status.NewStatus(err.Error(), http.StatusInternalServerError, status.WithDetails(status.Details{
							status.TableIndex:   tableIndex,
							status.RowNumber:    rowNum,
							status.ColumnNumber: cellNum,
						}))
						return false
					}
				}
			}
			// loop through the spans and populate table columns
			for i := 0; i < rowSpan; i++ {
				for j := 0; j < colSpan; j++ {
					row := rowNum + i
					nextAvailableCell := 0

					if _, ok := td[row]; !ok {
						td[row] = make(map[int]string)
					}

					columns := td[row]
					// check if column already has a value from a previous rowspan so we don't overrwite it
					// loop until we get an availalbe column
					// https://en.wikipedia.org/wiki/Help:Table#Combined_use_of_COLSPAN_and_ROWSPAN
					for columns[cellNum+j+nextAvailableCell] != "" {
						nextAvailableCell++
					}
					columns[cellNum+j+nextAvailableCell] = parseText(s.Text())
				}
			}
			return true
		})
		return err == nil
	})
	if err != nil {
		return nil, errorStatus
	}
	return td, nil
}

func parseText(s string) string {
	return strings.TrimSpace(s)
}

func getSpan(values []string) (int, error) {
	var err error
	var span int
	for _, v := range values {
		span, err = strconv.Atoi(v)
		if err == nil {
			return span, nil
		}
	}
	return 0, fmt.Errorf("no integer value in span attribute: %v", values)
}

package server

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable2json/internal/status"
	"golang.org/x/sync/errgroup"
)

type parseOptions struct {
	tables  []int
	keyrows int
}

func parse(ctx context.Context, tableSelection *goquery.Selection, input parseOptions) (interface{}, error) {
	var eg errgroup.Group
	switch len(input.tables) {
	case 0:
		resp := make([]interface{}, tableSelection.Length())
		tableSelection.Each(func(i int, selection *goquery.Selection) {
			eg.Go(func() error {
				td, err := parseTable(selection, i, input)
				if err != nil {
					return err
				}

				var tmp interface{}
				if input.keyrows >= 1 {
					tmp, err = formatKeyValue(td, input.keyrows, i)
					if err != nil {
						return err
					}
				} else {
					tmp = formatMatrix(td)
				}

				resp[i] = tmp
				return nil
			})
		})
		err := eg.Wait()
		if err != nil {
			return nil, err
		}
		return resp, nil
	default:
		resp := make([]interface{}, len(input.tables))
		for i, tableIndex := range input.tables {
			i := i
			tableIndex := tableIndex
			eg.Go(func() error {
				td, err := parseTable(tableSelection.Eq(tableIndex), tableIndex, input)
				if err != nil {
					return err
				}

				var tmp interface{}
				if input.keyrows >= 1 {
					tmp, err = formatKeyValue(td, input.keyrows, i)
					if err != nil {
						return err
					}
				} else {
					tmp = formatMatrix(td)
				}

				resp[i] = tmp
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

type cell struct {
	set   bool
	value string
}

func parseTable(tableSelection *goquery.Selection, tableIndex int, input parseOptions) (verbose, error) {
	td := make(verbose)

	errorStatus := status.Status{}
	var err error
	// for each row in the table
	tableSelection.Find("tr").EachWithBreak(func(rowNum int, s *goquery.Selection) bool {
		// find all th and td elements in the row
		var col int
		s.Find("th, td").EachWithBreak(func(cellNum int, s *goquery.Selection) bool {
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

			startCol := col

			// loop through the spans and populate table columns
			for i := 0; i < rowSpan; i++ {
				for j := 0; j < colSpan; j++ {
					row := rowNum + i
					nextAvailableCell := 0

					if _, ok := td[row]; !ok {
						td[row] = make(map[int]cell)
					}

					columns := td[row]

					// check if column already is already set from a previous rowspan so we don't overrwite it
					// loop until we get an availalbe column
					// https://en.wikipedia.org/wiki/Help:Table#Combined_use_of_COLSPAN_and_ROWSPAN
					for columns[startCol+j+nextAvailableCell].set {
						nextAvailableCell++
						if i == 0 {
							col++
						}
					}
					columns[startCol+j+nextAvailableCell] = cell{set: true, value: parseText(s)}
					if i == 0 {
						col++
					}
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

func parseText(s *goquery.Selection) string {
	return strings.TrimSpace(s.Text())
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

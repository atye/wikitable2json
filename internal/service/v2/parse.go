package v2

import (
	"context"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

type Tables struct {
	Caption string  `json:"caption"`
	Tables  []Table `json:"tables"`
}

type Table struct {
	Caption string     `json:"caption"`
	Data    [][]string `json:"data"`
}

type parseTableError struct {
	err        error
	tableIndex int
	rowNum     int
	cellNum    int
}

func (e *parseTableError) Error() string {
	return e.err.Error()
}

func parseTables(ctx context.Context, wikiTableSelection *goquery.Selection, tableIndices []string) (*Tables, error) {
	var eg errgroup.Group
	tables := new(Tables)
	switch len(tableIndices) {
	case 0:
		wikiTableSelection.Each(func(i int, selection *goquery.Selection) {
			eg.Go(func() error {
				table, err := parseTable(selection, i)
				if err != nil {
					return err
				}
				tables.Tables[i] = *table
				return nil
			})
		})
		err := eg.Wait()
		if err != nil {
			return nil, err
		}
		return tables, nil
	default:
		for i, tableIndex := range tableIndices {
			i := i
			tableIndex, err := strconv.Atoi(tableIndex)
			if err != nil {
				return nil, generalErr(err, http.StatusBadRequest)
			}
			eg.Go(func() error {
				table, err := parseTable(wikiTableSelection.Eq(tableIndex), tableIndex)
				if err != nil {
					return err
				}
				tables.Tables[i] = *table
				return nil
			})
		}
		err := eg.Wait()
		if err != nil {
			return nil, err
		}
		return tables, nil
	}
}

func parseTable(tableSelection *goquery.Selection, tableIndex int) (*Table, error) {
	table := new(Table)
	table.Caption = tableSelection.Find("caption").Text()
	ptErr := &parseTableError{}
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
					ptErr.err = err
					ptErr.rowNum = rowNum
					ptErr.cellNum = cellNum
					ptErr.tableIndex = tableIndex
					return false
				}
			}
			if attr := s.AttrOr("colspan", ""); attr != "" {
				colSpan, err = strconv.Atoi(attr)
				if err != nil {
					ptErr.err = err
					ptErr.rowNum = rowNum
					ptErr.cellNum = cellNum
					ptErr.tableIndex = tableIndex
					return false
				}
			}
			// loop through the spans and populate table columns
			for i := 0; i < rowSpan; i++ {
				for j := 0; j < colSpan; j++ {
					row := rowNum + i
					nextAvailableCell := 0

					columns := table.Data[row]
					for columns[cellNum+j+nextAvailableCell] != "" {
						nextAvailableCell++
					}
					columns[cellNum+j+nextAvailableCell] = parseText(s.Text())

					/*row := rowNum + i
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
					columns[int64(cellNum+j+nextAvailableCell)] = parseText(s.Text())*/
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
		return nil, tableParseErr(ptErr)
	}
	return nil, nil
}

func parseText(s string) string {
	return strings.TrimSpace(s)
}

package service

import (
	"context"
	"net/http"
	"strconv"
	"strings"
	"sync"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/sync/errgroup"
)

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

func parseTables(ctx context.Context, wikiTableSelection *goquery.Selection, tableIndices []string) ([]Table, error) {
	var eg errgroup.Group
	switch len(tableIndices) {
	case 0:
		tables := make([]Table, wikiTableSelection.Length())
		wikiTableSelection.Each(func(i int, selection *goquery.Selection) {
			eg.Go(func() error {
				table, err := parseTable(selection, i)
				if err != nil {
					return err
				}
				tables[i] = *table
				return nil
			})
		})
		err := eg.Wait()
		if err != nil {
			return nil, err
		}
		return tables, nil
	default:
		tables := make([]Table, len(tableIndices))
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
				tables[i] = *table
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
	table := &Table{}
	table.Caption = tableSelection.Find("caption").Text()

	dataMap := make(map[int]map[int]string)

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
				rowSpanTexts := strings.Split(attr, " ")
				if len(rowSpanTexts) > 0 {
					rowSpan, err = strconv.Atoi(rowSpanTexts[0])
					if err != nil {
						ptErr.err = err
						ptErr.rowNum = rowNum
						ptErr.cellNum = cellNum
						ptErr.tableIndex = tableIndex
						return false
					}
				}
			}
			if attr := s.AttrOr("colspan", ""); attr != "" {
				colSpanTexts := strings.Split(attr, " ")
				if len(colSpanTexts) > 0 {
					colSpan, err = strconv.Atoi(colSpanTexts[0])
					if err != nil {
						ptErr.err = err
						ptErr.rowNum = rowNum
						ptErr.cellNum = cellNum
						ptErr.tableIndex = tableIndex
						return false
					}
				}
			}
			// loop through the spans and populate table columns
			for i := 0; i < rowSpan; i++ {
				for j := 0; j < colSpan; j++ {
					row := rowNum + i
					nextAvailableCell := 0

					if _, ok := dataMap[row]; !ok {
						dataMap[row] = make(map[int]string)
					}

					columns := dataMap[row]
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
		return nil, tableParseErr(ptErr)
	}

	dataMapToData(dataMap, table)
	return table, nil
}

func parseText(s string) string {
	return strings.TrimSpace(s)
}

func dataMapToData(dataMap map[int]map[int]string, table *Table) {
	table.Data = make([][]string, len(dataMap))
	var wg sync.WaitGroup
	for i := 0; i < len(dataMap); i++ {
		wg.Add(1)
		go func(i int) {
			defer wg.Done()
			row := dataMap[i]
			table.Data[i] = make([]string, len(row))
			for j := 0; j < len(row); j++ {
				table.Data[i][j] = row[j]
			}
		}(i)
	}
	wg.Wait()
}

package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable2json/internal/api"
	"github.com/atye/wikitable2json/internal/cache"
	"github.com/atye/wikitable2json/internal/status"
	"golang.org/x/sync/errgroup"
)

var (
	classes = []string{
		"table.wikitable",
		"table.standard",
		"table.toccolours",
	}
)

type TableGetter interface {
	GetTablesMatrix(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([][][]string, error)
	GetTablesKeyValue(ctx context.Context, page string, lang string, cleanRef bool, keyRows int, tables ...int) ([][]map[string]string, error)
	SetUserAgent(string)
}

type Option func(*client)

func WithCache(capacity int, itemExpiration time.Duration, purgeEvery time.Duration) Option {
	return func(c *client) {
		c.cache = cache.NewCache(capacity, itemExpiration, purgeEvery)
	}
}

func WithHTTPClient(c *http.Client) Option {
	return func(tg *client) {
		tg.wikiAPI = api.NewWikiClient(api.BaseURL, api.WithHTTPClient(c))
	}
}

type wikiAPI interface {
	GetPageBytes(ctx context.Context, page, lang, userAgent string) ([]byte, error)
}

type client struct {
	wikiAPI   wikiAPI
	userAgent string
	cache     *cache.Cache
}

func NewTableGetter(userAgent string, options ...Option) TableGetter {
	c := &client{
		wikiAPI:   api.NewWikiClient(api.BaseURL),
		userAgent: userAgent,
	}

	for _, o := range options {
		o(c)
	}
	return c
}

func (c *client) GetTablesMatrix(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([][][]string, error) {
	tableSelection, err := c.getTableSelection(ctx, page, lang)
	if err != nil {
		return nil, handleErr(err)
	}

	if cleanRef {
		cleanReferences(tableSelection)
	}

	matrix, err := parse(tableSelection, 0, tables...)
	if err != nil {
		return nil, handleErr(err)
	}

	ret := [][][]string{}
	for _, v := range matrix {
		if m, ok := v.([][]string); ok {
			ret = append(ret, m)
		} else {
			return nil, status.NewStatus(fmt.Sprintf("unexpected return type %T", m), http.StatusInternalServerError)
		}
	}

	return ret, nil
}

func (c *client) GetTablesKeyValue(ctx context.Context, page string, lang string, cleanRef bool, keyRows int, tables ...int) ([][]map[string]string, error) {
	if keyRows < 1 {
		return nil, status.NewStatus("keyRows must be at least 1", http.StatusBadRequest)
	}

	tableSelection, err := c.getTableSelection(ctx, page, lang)
	if err != nil {
		return nil, handleErr(err)
	}

	if cleanRef {
		cleanReferences(tableSelection)
	}

	keyValue, err := parse(tableSelection, keyRows, tables...)
	if err != nil {
		return nil, handleErr(err)
	}

	ret := [][]map[string]string{}
	for _, v := range keyValue {
		if k, ok := v.([]map[string]string); ok {
			ret = append(ret, k)
		} else {
			return nil, status.NewStatus(fmt.Sprintf("unexpected return type %T", k), http.StatusInternalServerError)
		}
	}

	return ret, nil
}

func (c *client) SetUserAgent(agent string) {
	c.userAgent = agent
}

func (c *client) getTableSelection(ctx context.Context, page string, lang string) (*goquery.Selection, error) {
	var tableSelection *goquery.Selection
	var err error

	if c.cache != nil {
		var ok bool
		tableSelection, ok = c.cache.Get(page)
		if !ok {
			tableSelection, err = c.getTableSelectionFromAPI(ctx, page, lang)
			if err != nil {
				return nil, err
			}
			c.cache.Set(page, tableSelection)
		}
	} else {
		tableSelection, err = c.getTableSelectionFromAPI(ctx, page, lang)
		if err != nil {
			return nil, err
		}
	}

	return tableSelection, nil
}

func (c *client) getTableSelectionFromAPI(ctx context.Context, page string, lang string) (*goquery.Selection, error) {
	b, err := c.wikiAPI.GetPageBytes(ctx, page, lang, c.userAgent)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	doc.Find(".mw-empty-elt").Remove()
	doc.Find("style").Remove()

	return doc.Find(strings.Join(classes, ", ")), nil
}

func cleanReferences(tables *goquery.Selection) {
	tables.Find(".reference").Remove()

	tables.Find("sup").Each(func(_ int, s *goquery.Selection) {
		s.Find("a").EachWithBreak(func(_ int, anchor *goquery.Selection) bool {
			if v, ok := anchor.Attr("title"); ok {
				if v == "Wikipedia:Citation needed" {
					s.Remove()
					return false
				}
			}
			return true
		})
	})
}

func parse(tableSelection *goquery.Selection, keyRows int, tables ...int) ([]interface{}, error) {
	var ret []interface{}

	var eg errgroup.Group
	switch len(tables) {
	case 0:
		ret = make([]interface{}, tableSelection.Length())
		tableSelection.Each(func(i int, selection *goquery.Selection) {
			eg.Go(func() error {
				td, err := parseTable(selection, i)
				if err != nil {
					return err
				}

				var tmp interface{}
				if keyRows >= 1 {
					tmp, err = formatKeyValue(td, keyRows, i)
					if err != nil {
						return err
					}
				} else {
					tmp = formatMatrix(td)
				}

				ret[i] = tmp
				return nil
			})
		})
	default:
		ret = make([]interface{}, len(tables))
		for i, tableIndex := range tables {
			i := i
			tableIndex := tableIndex
			eg.Go(func() error {
				td, err := parseTable(tableSelection.Eq(tableIndex), tableIndex)
				if err != nil {
					return err
				}

				var tmp interface{}
				if keyRows >= 1 {
					tmp, err = formatKeyValue(td, keyRows, i)
					if err != nil {
						return err
					}
				} else {
					tmp = formatMatrix(td)
				}

				ret[i] = tmp
				return nil
			})
		}
	}

	err := eg.Wait()
	if err != nil {
		return nil, err
	}
	return ret, nil
}

type cell struct {
	set   bool
	value string
}

func parseTable(tableSelection *goquery.Selection, tableIndex int) (verbose, error) {
	td := make(verbose)

	tableClass := getTableClass(tableSelection)
	if tableClass == "" {
		return td, nil
	}

	errorStatus := status.Status{}
	var err error
	tableSelection.Find(fmt.Sprintf("table.%s > thead > tr, table.%s > tbody > tr", tableClass, tableClass)).EachWithBreak(func(rowNum int, row *goquery.Selection) bool {
		var col int
		if _, ok := td[rowNum]; !ok {
			td[rowNum] = make(map[int]cell)
		}

		row.Find(fmt.Sprintf("table.%s > thead > tr > th, table.%s > thead > tr > td, table.%s > tbody > tr > th, table.%s > tbody > tr > td", tableClass, tableClass, tableClass, tableClass)).EachWithBreak(func(cellNum int, s *goquery.Selection) bool {
			rowSpan := 1
			colSpan := 1
			if attr := s.AttrOr("rowspan", ""); attr != "" {
				rowSpanTexts := strings.Split(attr, " ")
				if len(rowSpanTexts) > 0 {
					rowSpan, err = getSpan(rowSpanTexts)
					if err != nil {
						errorStatus = status.NewStatus(err.Error(), http.StatusInternalServerError, status.WithDetails(status.Details{
							status.TableIndex:  tableIndex,
							status.RowIndex:    rowNum,
							status.ColumnIndex: cellNum,
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
							status.TableIndex:  tableIndex,
							status.RowIndex:    rowNum,
							status.ColumnIndex: cellNum,
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

func getTableClass(table *goquery.Selection) string {
	v, ok := table.Attr("class")
	if !ok {
		return ""
	}

	values := strings.Split(v, " ")
	for _, s := range values {
		if s == "wikitable" || s == "toccolours" || s == "standard" {
			return s
		}
	}
	return ""
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

func handleErr(err error) status.Status {
	var s status.Status
	if errors.As(err, &s) {
		return s
	} else {
		return status.NewStatus(err.Error(), http.StatusInternalServerError)
	}
}

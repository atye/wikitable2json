package client

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable2json/internal/server/api"
	"github.com/atye/wikitable2json/internal/server/status"
	"github.com/atye/wikitable2json/pkg/cache"
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
	GetTablesMatrix(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([]Matrix, error)
	GetTablesKeyValue(ctx context.Context, page string, lang string, cleanRef bool, keyRows int, tables ...int) ([]KeyValue, error)
	SetUserAgent(string)
}

type Option func(*client)

func WithCache(cache *cache.Cache) Option {
	return func(c *client) {
		c.cache = cache
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

func (c *client) GetTablesMatrix(ctx context.Context, page string, lang string, cleanRef bool, tables ...int) ([]Matrix, error) {
	tableSelection, err := c.getTableSelection(ctx, page, lang, cleanRef)
	if err != nil {
		return nil, err
	}

	matrix, err := parse(ctx, tableSelection, 0, tables...)
	if err != nil {
		return nil, err
	}

	var ret []Matrix
	for _, v := range matrix {
		if m, ok := v.(Matrix); ok {
			ret = append(ret, m)
		} else {
			return nil, fmt.Errorf("unexpected return type %T", m)
		}
	}

	return ret, nil
}

func (c *client) GetTablesKeyValue(ctx context.Context, page string, lang string, cleanRef bool, keyRows int, tables ...int) ([]KeyValue, error) {
	if keyRows < 1 {
		return nil, fmt.Errorf("keyRows must be at least 1")
	}

	tableSelection, err := c.getTableSelection(ctx, page, lang, cleanRef)
	if err != nil {
		return nil, err
	}

	keyValue, err := parse(ctx, tableSelection, keyRows, tables...)
	if err != nil {
		return nil, err
	}

	var ret []KeyValue
	for _, v := range keyValue {
		if k, ok := v.(KeyValue); ok {
			ret = append(ret, k)
		} else {
			return nil, fmt.Errorf("unexpected return type %T", k)
		}
	}

	return ret, nil
}

func (c *client) SetUserAgent(agent string) {
	c.userAgent = agent
}

func (c *client) getTableSelection(ctx context.Context, page string, lang string, cleanRef bool) (*goquery.Selection, error) {
	var tableSelection *goquery.Selection
	var err error

	if c.cache != nil {
		var ok bool
		tableSelection, ok = c.cache.Get(page)
		if !ok {
			tableSelection, err = c.getTableSelectionFromAPI(ctx, page, lang, cleanRef)
			if err != nil {
				return nil, err
			}
			c.cache.Set(page, tableSelection)
		}
	} else {
		tableSelection, err = c.getTableSelectionFromAPI(ctx, page, lang, cleanRef)
		if err != nil {
			return nil, err
		}
	}

	if cleanRef {
		cleanReferences(tableSelection)
	}

	return tableSelection, nil
}

func (c *client) getTableSelectionFromAPI(ctx context.Context, page string, lang string, cleanRef bool) (*goquery.Selection, error) {
	b, err := c.wikiAPI.GetPageBytes(ctx, page, lang, c.userAgent)
	if err != nil {
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
	if err != nil {
		return nil, err
	}
	doc.Find(".mw-empty-elt").Remove()

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

func parse(ctx context.Context, tableSelection *goquery.Selection, keyRows int, tables ...int) ([]interface{}, error) {
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

				var tmp interface{}
				if keyRows >= 1 {
					tmp, err = formatKeyValue(td, keyRows, i)
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
		resp := make([]interface{}, len(tables))
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

func parseTable(tableSelection *goquery.Selection, tableIndex int) (verbose, error) {
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

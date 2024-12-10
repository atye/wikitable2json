package client

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/atye/wikitable2json/pkg/client/status"
	"golang.org/x/net/html"
	"golang.org/x/sync/errgroup"
)

var (
	classes = []string{
		"table.wikitable",
		"table.standard",
		"table.toccolours",
	}

	defaultUserAgent = "github.com/atye/wikitable2json"

	apiURL      = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
	getApiURLFn = getApiURL
)

type Client struct {
	http      *http.Client
	userAgent string
}

type ClientOption func(*Client)

func WithHTTPClient(c *http.Client) ClientOption {
	return func(tg *Client) {
		tg.http = c
	}
}

type tableOptions struct {
	cleanRef  bool
	brNewLine bool
	tables    []int
}

type TableOption func(*tableOptions)

func WithCleanReferences() TableOption {
	return func(to *tableOptions) {
		to.cleanRef = true
	}
}

func WithBRNewLine() TableOption {
	return func(to *tableOptions) {
		to.brNewLine = true
	}
}

func WithTables(tables ...int) TableOption {
	return func(to *tableOptions) {
		to.tables = tables
	}
}

func NewClient(userAgent string, options ...ClientOption) *Client {
	c := &Client{
		http:      http.DefaultClient,
		userAgent: userAgent,
	}

	for _, o := range options {
		o(c)
	}
	return c
}

func (c *Client) GetMatrix(ctx context.Context, page string, lang string, options ...TableOption) ([][][]string, error) {
	tableSelection, err := c.getTableSelectionFromAPI(ctx, page, lang)
	if err != nil {
		return nil, handleErr(err)
	}

	to := new(tableOptions)
	for _, o := range options {
		o(to)
	}

	if to.cleanRef {
		cleanReferences(tableSelection)
	}

	matrix, err := parse(tableSelection, 0, false, to.brNewLine, to.tables...)
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

func (c *Client) GetMatrixVerbose(ctx context.Context, page string, lang string, options ...TableOption) ([][][]Verbose, error) {
	tableSelection, err := c.getTableSelectionFromAPI(ctx, page, lang)
	if err != nil {
		return nil, handleErr(err)
	}

	to := new(tableOptions)
	for _, o := range options {
		o(to)
	}

	if to.cleanRef {
		cleanReferences(tableSelection)
	}

	matrix, err := parse(tableSelection, 0, true, to.brNewLine, to.tables...)
	if err != nil {
		return nil, handleErr(err)
	}

	ret := [][][]Verbose{}
	for _, v := range matrix {
		if m, ok := v.([][]Verbose); ok {
			ret = append(ret, m)
		} else {
			return nil, status.NewStatus(fmt.Sprintf("unexpected return type %T", m), http.StatusInternalServerError)
		}
	}

	return ret, nil
}

func (c *Client) GetKeyValue(ctx context.Context, page string, lang string, keyRows int, options ...TableOption) ([][]map[string]string, error) {
	if keyRows < 1 {
		return nil, status.NewStatus("keyRows must be at least 1", http.StatusBadRequest)
	}

	tableSelection, err := c.getTableSelectionFromAPI(ctx, page, lang)
	if err != nil {
		return nil, handleErr(err)
	}

	to := new(tableOptions)
	for _, o := range options {
		o(to)
	}

	if to.cleanRef {
		cleanReferences(tableSelection)
	}

	keyValue, err := parse(tableSelection, keyRows, false, to.brNewLine, to.tables...)
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

func (c *Client) GetKeyValueVerbose(ctx context.Context, page string, lang string, keyRows int, options ...TableOption) ([][]map[string]Verbose, error) {
	if keyRows < 1 {
		return nil, status.NewStatus("keyRows must be at least 1", http.StatusBadRequest)
	}

	tableSelection, err := c.getTableSelectionFromAPI(ctx, page, lang)
	if err != nil {
		return nil, handleErr(err)
	}

	to := new(tableOptions)
	for _, o := range options {
		o(to)
	}

	if to.cleanRef {
		cleanReferences(tableSelection)
	}

	keyValue, err := parse(tableSelection, keyRows, true, to.brNewLine, to.tables...)
	if err != nil {
		return nil, handleErr(err)
	}

	ret := [][]map[string]Verbose{}
	for _, v := range keyValue {
		if k, ok := v.([]map[string]Verbose); ok {
			ret = append(ret, k)
		} else {
			return nil, status.NewStatus(fmt.Sprintf("unexpected return type %T", k), http.StatusInternalServerError)
		}
	}

	return ret, nil
}

func (c *Client) SetUserAgent(userAgent string) {
	c.userAgent = userAgent
}

func (c *Client) getTableSelectionFromAPI(ctx context.Context, page string, lang string) (*goquery.Selection, error) {
	b, err := c.getPageData(ctx, page, lang)
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

func (c *Client) getPageData(ctx context.Context, page string, lang string) ([]byte, error) {
	u, err := url.Parse(getApiURLFn(lang, url.QueryEscape(page)))
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	agent := defaultUserAgent
	if c.userAgent != "" {
		agent = c.userAgent
	}
	req.Header.Add("User-Agent", agent)

	resp, err := c.http.Do(req)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError, status.WithDetails(status.Details{
			status.Page: page,
		}))
	}
	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError, status.WithDetails(status.Details{
			status.Page: page,
		}))
	}

	if resp.StatusCode != http.StatusOK {
		return nil, status.NewStatus(string(b), resp.StatusCode, status.WithDetails(status.Details{
			status.Page: page,
		}))
	}

	return b, nil
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

func parse(tableSelection *goquery.Selection, keyRows int, verbose bool, brNewLine bool, tables ...int) ([]interface{}, error) {
	var ret []interface{}

	var eg errgroup.Group
	switch len(tables) {
	case 0:
		ret = make([]interface{}, tableSelection.Length())
		tableSelection.Each(func(i int, selection *goquery.Selection) {
			eg.Go(func() error {
				td, err := parseTable(selection, i, brNewLine)
				if err != nil {
					return err
				}

				tmp, err := formatParsedTable(td, verbose, keyRows, i)
				if err != nil {
					return err
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
				td, err := parseTable(tableSelection.Eq(tableIndex), tableIndex, brNewLine)
				if err != nil {
					return err
				}
				tmp, err := formatParsedTable(td, verbose, keyRows, i)
				if err != nil {
					return err
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

func formatParsedTable(data parsed, verbose bool, keyRows int, tableIndex int) (interface{}, error) {
	var tmp interface{}
	var err error
	if keyRows >= 1 {
		if verbose {
			tmp, err = formatKeyValueVerbose(data, keyRows, tableIndex)
			if err != nil {
				return nil, err
			}
		} else {
			tmp, err = formatKeyValue(data, keyRows, tableIndex)
			if err != nil {
				return nil, err
			}
		}
	} else {
		if verbose {
			tmp = formatMatrixVerbose(data)
		} else {
			tmp = formatMatrix(data)
		}
	}
	return tmp, nil
}

type cell struct {
	set   bool
	text  string
	links []Link
}

func parseTable(tableSelection *goquery.Selection, tableIndex int, brIsNewLine bool) (parsed, error) {
	td := make(parsed)

	tableClass := getTableClass(tableSelection)
	if tableClass == "" {
		return td, nil
	}

	parseNonTextNodeFuncs := []func(*html.Node) string{}
	if brIsNewLine {
		parseNonTextNodeFuncs = append(parseNonTextNodeFuncs, brNewLine)
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
					columns[startCol+j+nextAvailableCell] = cell{
						set:   true,
						text:  parseText(s, parseNonTextNodeFuncs...),
						links: parseLink(s, parseNonTextNodeFuncs...)}
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

func parseText(s *goquery.Selection, parseNonTextNode ...func(*html.Node) string) string {
	var buf bytes.Buffer

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.TextNode {
			buf.WriteString(n.Data)
		} else if len(parseNonTextNode) > 0 {
			for _, parseFn := range parseNonTextNode {
				buf.WriteString(parseFn(n))
			}
		}
		if n.FirstChild != nil {
			for c := n.FirstChild; c != nil; c = c.NextSibling {
				f(c)
			}
		}
	}
	for _, n := range s.Nodes {
		f(n)
	}

	return buf.String()
}

func parseLink(s *goquery.Selection, parseNonTextNode ...func(*html.Node) string) []Link {
	var ret []Link
	s.Find("a").Each(func(_ int, anchor *goquery.Selection) {
		if v, ok := anchor.Attr("href"); ok {
			if v != "" {
				ret = append(ret, Link{Href: v, Text: parseText(anchor, parseNonTextNode...)})
			}
		}
	})
	return ret
}

func brNewLine(n *html.Node) string {
	if n.Data == "br" {
		return "\n"
	}
	return ""
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
	}
	return status.NewStatus(err.Error(), http.StatusInternalServerError)
}

func getApiURL(lang, page string) string {
	return fmt.Sprintf(apiURL, page)
}

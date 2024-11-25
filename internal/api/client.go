package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"

	"github.com/atye/wikitable2json/internal/status"
)

var (
	BaseURL = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
)

var GetEndpoint = func(lang, page string) string {
	return fmt.Sprintf(BaseURL, lang, url.QueryEscape(page))
}

type WikiClient struct {
	client *http.Client
}

type Option func(c *WikiClient)

func WithHTTPClient(c *http.Client) Option {
	return func(wc *WikiClient) {
		wc.client = c
	}
}

func NewWikiClient(options ...Option) *WikiClient {
	wc := &WikiClient{
		client: &http.Client{Timeout: 10 * time.Second},
	}

	for _, o := range options {
		o(wc)
	}
	return wc
}

func (c *WikiClient) GetPage(ctx context.Context, page, lang, userAgent string) ([]byte, error) {
	u, err := url.Parse(GetEndpoint(lang, page))
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	agent := "github.com/atye/wikitable2json"
	if userAgent != "" {
		agent = userAgent
	}
	req.Header.Add("User-Agent", agent)

	resp, err := c.client.Do(req)
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

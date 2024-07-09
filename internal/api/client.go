package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/atye/wikitable2json/internal/status"
)

var (
	BaseURL = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
)

type WikiClient struct {
	client   *http.Client
	endpoint string
}

type Option func(c *WikiClient)

func WithHTTPClient(c *http.Client) Option {
	return func(wc *WikiClient) {
		wc.client = c
	}
}

func NewWikiClient(endpoint string, options ...Option) *WikiClient {
	wc := &WikiClient{
		client:   &http.Client{Timeout: 10 * time.Second},
		endpoint: endpoint,
	}

	for _, o := range options {
		o(wc)
	}
	return wc
}

func (c *WikiClient) GetPageBytes(ctx context.Context, page, lang, userAgent string) ([]byte, error) {
	addr, err := buildURL(c.endpoint, page, lang)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
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

func buildURL(endpoint, page, lang string) (string, error) {
	if strings.Contains(endpoint, "wikipedia.org") {
		u, err := url.Parse(fmt.Sprintf(BaseURL, lang, url.QueryEscape(page)))
		if err != nil {
			return "", err
		}
		return u.String(), nil
	}

	u, err := url.Parse(fmt.Sprintf("%s/api/rest_v1/page/html/%s", endpoint, url.QueryEscape(page)))
	if err != nil {
		return "", err
	}
	return u.String(), nil
}

package data

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/atye/wikitable-api/internal/status"
)

const (
	defaultUserAgent = "Go-http-client/1.1"
)

var (
	BaseURL = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
)

type WikiClient struct {
	client   *http.Client
	endpoint string
}

func NewWikiClient(endpoint string) WikiClient {
	return WikiClient{
		client:   http.DefaultClient,
		endpoint: endpoint,
	}
}

func (c WikiClient) GetPageData(ctx context.Context, page, lang, userAgent string) (io.ReadCloser, error) {
	addr, err := buildURL(c.endpoint, page, lang)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, addr, nil)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	if userAgent != defaultUserAgent {
		req.Header.Add("User-Agent", userAgent)
	}

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}

	if resp.StatusCode != http.StatusOK {
		defer resp.Body.Close()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, status.NewStatus(err.Error(), resp.StatusCode, status.WithDetails(status.Details{
				status.Page: page,
			}))
		}
		return nil, status.NewStatus(string(body), resp.StatusCode, status.WithDetails(status.Details{
			status.Page: page,
		}))
	}

	return resp.Body, nil
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

package api

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/atye/wikitable2json/internal/server/status"
)

type WikiClient struct {
	client   *http.Client
	endpoint string
}

var (
	BaseURL = "https://%s.wikipedia.org/api/rest_v1/page/html/%s"
)

func NewWikiClient(endpoint string) WikiClient {
	return WikiClient{
		client:   http.DefaultClient,
		endpoint: endpoint,
	}
}

func (c WikiClient) GetPageBytes(ctx context.Context, page, lang, userAgent string) ([]byte, error) {
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
		return nil, status.NewStatus(err.Error(), http.StatusInternalServerError)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
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

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
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

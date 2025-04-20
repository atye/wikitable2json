package metrics

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
)

var (
	apiURL = "https://www.google-analytics.com/mp/collect?measurement_id=%s&api_secret=%s"
)

type GoogleClient struct {
	measurementID string
	apiSecret     string
	httpClient    *http.Client
}

type gaEvent struct {
	ClientID string  `json:"client_id"`
	Events   []event `json:"events"`
}

type event struct {
	Name   string                 `json:"name"`
	Params map[string]interface{} `json:"params"`
}

func NewGoogleClient(measurementID, apiSecret string, client *http.Client) *GoogleClient {
	return &GoogleClient{
		measurementID: measurementID,
		apiSecret:     apiSecret,
		httpClient:    client,
	}
}

func (c *GoogleClient) Publish(code int, ip, page, lang string) error {
	hash := sha256.Sum256([]byte(ip))

	event := gaEvent{
		ClientID: hex.EncodeToString(hash[:]),
		Events: []event{
			{
				Name: "page_request",
				Params: map[string]interface{}{
					"page": page,
					"lang": lang,
					"code": code,
				},
			},
		},
	}

	body, err := json.Marshal(event)
	if err != nil {
		return err
	}

	// https://developers.google.com/analytics/devguides/collection/protocol/ga4/reference#response_codes
	// The Measurement Protocol always returns a 2xx status code if the HTTP request was received.
	// The Measurement Protocol does not return an error code if the payload data was malformed, or if the data in the payload was incorrect or was not processed by Google Analytics.
	_, err = c.httpClient.Post(fmt.Sprintf(apiURL, c.measurementID, c.apiSecret), "application/json", bytes.NewBuffer(body))
	return err
}

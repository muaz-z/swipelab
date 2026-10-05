package network

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

type Client struct {
	baseURL    string
	httpClient *http.Client
}

func NewClient(baseURL string) *Client {
	return &Client{
		baseURL:    baseURL,
		httpClient: &http.Client{},
	}
}

func (c *Client) Authorize(
	ctx context.Context,
	request AuthorizationRequest,
) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal authorization request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.baseURL+"/authorizations",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create authorization request: %w", err)

	}

	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("send authorization request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("card network returned status %d", response.StatusCode)
	}

	return nil
}

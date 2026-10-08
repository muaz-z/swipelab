package issuer

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

type Client struct {
	httpClient *http.Client
}

func NewClient() *Client {
	return &Client{
		httpClient: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *Client) Authorize(
	ctx context.Context,
	issuerURL string,
	request any,
) error {
	body, err := json.Marshal(request)
	if err != nil {
		return fmt.Errorf("marshal issuer request: %w", err)
	}

	httpRequest, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		issuerURL+"/authorizations",
		bytes.NewReader(body),
	)
	if err != nil {
		return fmt.Errorf("create issuer request: %w", err)
	}

	httpRequest.Header.Set("Content-Type", "application/json")

	response, err := c.httpClient.Do(httpRequest)
	if err != nil {
		return fmt.Errorf("send issuer request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode < 200 || response.StatusCode >= 300 {
		return fmt.Errorf("issuer returned HTTP %d", response.StatusCode)
	}

	return nil
}

package bento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

// GetBroadcasts retrieves all broadcasts
func (c *Client) GetBroadcasts(ctx context.Context) ([]BroadcastData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/fetch/broadcasts", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result struct {
		Broadcasts []BroadcastData `json:"broadcasts"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result.Broadcasts, nil
}

// CreateBroadcast creates a new broadcast
func (c *Client) CreateBroadcast(ctx context.Context, broadcasts []BroadcastData) error {
	if len(broadcasts) == 0 {
		return ErrInvalidRequest
	}

	// Validate broadcasts before sending
	for _, broadcast := range broadcasts {
		if broadcast.Name == "" {
			return fmt.Errorf("%w: broadcast name is required", ErrInvalidRequest)
		}
		if broadcast.Subject == "" {
			return fmt.Errorf("%w: broadcast subject is required", ErrInvalidRequest)
		}
		if broadcast.Content == "" {
			return fmt.Errorf("%w: broadcast content is required", ErrInvalidRequest)
		}
		if _, err := mail.ParseAddress(broadcast.From.Email); err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidEmail, broadcast.From.Email)
		}
		if broadcast.BatchSizePerHour <= 0 {
			return fmt.Errorf("%w: batch size must be positive", ErrInvalidBatchSize)
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"broadcasts": broadcasts,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/batch/broadcasts", c.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	return nil
}

package bento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

// TrackEvent sends tracking events to Bento
func (c *Client) TrackEvent(ctx context.Context, events []EventData) error {
	if len(events) == 0 {
		return ErrInvalidRequest
	}

	// Validate all emails in events before sending
	for _, event := range events {
		if _, err := mail.ParseAddress(event.Email); err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidEmail, event.Email)
		}
		if event.Type == "" {
			return fmt.Errorf("%w: event type is required", ErrInvalidRequest)
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"events": events,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/batch/events", c.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result struct {
		Results int `json:"results"`
		Failed  int `json:"failed"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Failed > 0 {
		return fmt.Errorf("event tracking partially failed: %d succeeded, %d failed", result.Results, result.Failed)
	}

	return nil
}

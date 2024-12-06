package bento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

// SubscriberInput represents the data structure for creating/importing subscribers
type SubscriberInput struct {
	Email      string                 `json:"email"`
	FirstName  string                 `json:"first_name,omitempty"`
	LastName   string                 `json:"last_name,omitempty"`
	Tags       []string               `json:"tags,omitempty"`
	RemoveTags []string               `json:"remove_tags,omitempty"`
	Fields     map[string]interface{} `json:"fields,omitempty"`
}

// FindSubscriber retrieves a subscriber by email
func (c *Client) FindSubscriber(ctx context.Context, email string) (*SubscriberData, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEmail, email)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/fetch/subscribers", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("email", email)
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var response struct {
		Data SubscriberData `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	if response.Data.ID == "" {
		return nil, fmt.Errorf("subscriber not found: %s", email)
	}

	return &response.Data, nil
}

// CreateSubscriber creates a new subscriber
func (c *Client) CreateSubscriber(ctx context.Context, input *SubscriberInput) (*SubscriberData, error) {
	if _, err := mail.ParseAddress(input.Email); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEmail, input.Email)
	}

	body, err := json.Marshal(map[string]interface{}{
		"subscriber": input,
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/fetch/subscribers", c.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var response struct {
		Data SubscriberData `json:"data"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &response.Data, nil
}

// ImportSubscribers imports multiple subscribers in batch
func (c *Client) ImportSubscribers(ctx context.Context, subscribers []*SubscriberInput) error {
	if len(subscribers) == 0 {
		return ErrInvalidRequest
	}

	// Validate all emails before sending
	for _, sub := range subscribers {
		if _, err := mail.ParseAddress(sub.Email); err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidEmail, sub.Email)
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"subscribers": subscribers,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/batch/subscribers", c.baseURL), bytes.NewBuffer(body))
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
		return fmt.Errorf("import partially failed: %d succeeded, %d failed", result.Results, result.Failed)
	}

	return nil
}

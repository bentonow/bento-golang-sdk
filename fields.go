package bento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetFields retrieves all custom fields
func (c *Client) GetFields(ctx context.Context) ([]FieldData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/fetch/fields", c.baseURL), nil)
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

	var result FieldsResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}

// CreateField creates a new custom field
func (c *Client) CreateField(ctx context.Context, key string) (*FieldData, error) {
	if key == "" {
		return nil, fmt.Errorf("%w: field key is required", ErrInvalidRequest)
	}

	body, err := json.Marshal(map[string]interface{}{
		"field": map[string]string{
			"key": key,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/fetch/fields", c.baseURL), bytes.NewBuffer(body))
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

	var result struct {
		Data FieldData `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result.Data, nil
}

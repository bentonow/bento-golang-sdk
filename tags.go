package bento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetTags retrieves all tags
func (c *Client) GetTags(ctx context.Context) ([]TagData, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/fetch/tags", c.baseURL), nil)
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
		Data []TagData `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return result.Data, nil
}

// CreateTag creates a new tag
func (c *Client) CreateTag(ctx context.Context, tagName string) (*TagData, error) {
	if tagName == "" {
		return nil, fmt.Errorf("%w: tag name is required", ErrInvalidRequest)
	}

	body, err := json.Marshal(map[string]interface{}{
		"tag": map[string]string{
			"name": tagName,
		},
	})
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/fetch/tags", c.baseURL), bytes.NewBuffer(body))
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
		Data TagData `json:"data"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result.Data, nil
}

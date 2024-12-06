package bento

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"net/http"
	"net/mail"
)

// GetBlacklistStatus checks domain or IP address blacklist status
func (c *Client) GetBlacklistStatus(ctx context.Context, data *BlacklistData) (map[string]interface{}, error) {
	if data.Domain == "" && data.IPAddress == "" {
		return nil, fmt.Errorf("%w: either domain or IP address is required", ErrInvalidRequest)
	}

	if data.IPAddress != "" {
		if ip := net.ParseIP(data.IPAddress); ip == nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidIPAddress, data.IPAddress)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/experimental/blacklist.json", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	if data.Domain != "" {
		q.Add("domain", data.Domain)
	}
	if data.IPAddress != "" {
		q.Add("ip", data.IPAddress)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// ValidateEmail validates an email address
func (c *Client) ValidateEmail(ctx context.Context, data *ValidationData) (*ValidationResponse, error) {
	if _, err := mail.ParseAddress(data.EmailAddress); err != nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidEmail, data.EmailAddress)
	}

	if data.IPAddress != "" {
		if ip := net.ParseIP(data.IPAddress); ip == nil {
			return nil, fmt.Errorf("%w: %s", ErrInvalidIPAddress, data.IPAddress)
		}
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/experimental/validation", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("email", data.EmailAddress)
	if data.FullName != "" {
		q.Add("name", data.FullName)
	}
	if data.UserAgent != "" {
		q.Add("user_agent", data.UserAgent)
	}
	if data.IPAddress != "" {
		q.Add("ip", data.IPAddress)
	}
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result ValidationResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}

	return &result, nil
}

// GetContentModeration performs content moderation
func (c *Client) GetContentModeration(ctx context.Context, content string) (map[string]interface{}, error) {
	if content == "" {
		return nil, fmt.Errorf("%w: content is required", ErrInvalidContent)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/experimental/content_moderation", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("content", content)
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetGender predicts gender from a name
func (c *Client) GetGender(ctx context.Context, fullName string) (map[string]interface{}, error) {
	if fullName == "" {
		return nil, fmt.Errorf("%w: full name is required", ErrInvalidName)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/experimental/gender", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("name", fullName)
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GeoLocateIP performs IP geolocation
func (c *Client) GeoLocateIP(ctx context.Context, ipAddress string) (map[string]interface{}, error) {
	if ip := net.ParseIP(ipAddress); ip == nil {
		return nil, fmt.Errorf("%w: %s", ErrInvalidIPAddress, ipAddress)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/experimental/geolocation", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("ip", ipAddress)
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

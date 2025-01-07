package bento

import (
	"fmt"
	"net/http"
	"time"
)

// Client is the main entry point for the Bento SDK
type Client struct {
	baseURL    string
	httpClient HTTPDoer
	config     *Config
}

// HTTPDoer interface for HTTP client implementations
type HTTPDoer interface {
	Do(*http.Request) (*http.Response, error)
}

// Config holds the configuration for the Bento client
type Config struct {
	PublishableKey string
	SecretKey      string
	SiteUUID       string
	Timeout        time.Duration
}

// NewClient creates a new Bento client with the given configuration
func NewClient(config *Config) (*Client, error) {
	if config.PublishableKey == "" || config.SecretKey == "" || config.SiteUUID == "" {
		return nil, ErrInvalidConfig
	}

	// Validate timeout value
	if config.Timeout < 0 {
		return nil, fmt.Errorf("timeout must be non-negative")
	}

	// Set default timeout if none provided
	if config.Timeout == 0 {
		config.Timeout = 10 * time.Second
	}

	return &Client{
		baseURL: "https://app.bentonow.com/api/v1",
		httpClient: &http.Client{
			Timeout: config.Timeout,
		},
		config: config,
	}, nil
}

// do executes an HTTP request with proper context handling
func (c *Client) do(req *http.Request) (*http.Response, error) {
	// Check if context is already cancelled/timeout
	if err := req.Context().Err(); err != nil {
		return nil, err
	}

	req.SetBasicAuth(c.config.PublishableKey, c.config.SecretKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "bento-go-"+c.config.SiteUUID)

	q := req.URL.Query()
	q.Add("site_uuid", c.config.SiteUUID)
	req.URL.RawQuery = q.Encode()

	return c.httpClient.Do(req)
}

// SetHTTPClient sets a custom HTTP client
func (c *Client) SetHTTPClient(client HTTPDoer) error {
	if client == nil {
		return fmt.Errorf("HTTP client cannot be nil")
	}
	c.httpClient = client
	return nil
}

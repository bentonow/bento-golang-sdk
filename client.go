package bento

import (
	"net/http"
	"time"
)

// Client is the main entry point for the Bento SDK
type Client struct {
	baseURL    string
	httpClient *http.Client
	config     *Config
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

// do executes an HTTP request and returns the response
func (c *Client) do(req *http.Request) (*http.Response, error) {
	req.SetBasicAuth(c.config.PublishableKey, c.config.SecretKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("User-Agent", "bento-go-"+c.config.SiteUUID)

	q := req.URL.Query()
	q.Add("site_uuid", c.config.SiteUUID)
	req.URL.RawQuery = q.Encode()

	return c.httpClient.Do(req)
}

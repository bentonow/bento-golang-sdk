package bento

import (
	"fmt"
	"net/http"
	"time"
	"strings"
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
	var missingFields []string

	if config.PublishableKey == "" {
		missingFields = append(missingFields, "PublishableKey")
	}
	if config.SecretKey == "" {
		missingFields = append(missingFields, "SecretKey")
	}
	if config.SiteUUID == "" {
		missingFields = append(missingFields, "SiteUUID")
	}

	if len(missingFields) > 0 {
		return nil, fmt.Errorf("%w: %s", ErrInvalidConfig, strings.Join(missingFields, ", "))
	}

	if l := len(strings.Trim(config.PublishableKey, "\"")); l < 28 || l > 36 {
        return nil, fmt.Errorf("%w: PublishableKey must be between 28 and 36 characters (got %d)", ErrInvalidKeyLength, l)
    }
    if l := len(strings.Trim(config.SecretKey, "\"")); l < 28 || l > 36 {
        return nil, fmt.Errorf("%w: SecretKey must be between 28 and 36 characters (got %d)", ErrInvalidKeyLength, l)
    }
    if l := len(strings.Trim(config.SiteUUID, "\"")); l < 28 || l > 36 {
        return nil, fmt.Errorf("%w: SiteUUID must be between 28 and 36 characters (got %d)", ErrInvalidKeyLength, l)
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

    resp, err := c.httpClient.Do(req)
    if err != nil {
        return nil, fmt.Errorf("request failed: %w", err)
    }

    // Provide specific error messages based on status code
    switch resp.StatusCode {
    case http.StatusOK, http.StatusCreated:
        return resp, nil
    case http.StatusUnauthorized:
        return nil, fmt.Errorf("%w: invalid authentication credentials (401)", ErrAPIResponse)
    case http.StatusForbidden:
        return nil, fmt.Errorf("%w: access forbidden (403)", ErrAPIResponse)
    case http.StatusNotFound:
        return nil, fmt.Errorf("%w: resource not found (404)", ErrAPIResponse)
    case http.StatusBadRequest:
        return nil, fmt.Errorf("%w: invalid request parameters (400)", ErrAPIResponse)
    case http.StatusTooManyRequests:
        return nil, fmt.Errorf("%w: rate limit exceeded (429)", ErrAPIResponse)
    case http.StatusInternalServerError:
        return nil, fmt.Errorf("%w: server error (500)", ErrAPIResponse)
    case http.StatusServiceUnavailable:
        return nil, fmt.Errorf("%w: service unavailable (503)", ErrAPIResponse)
    default:
        return nil, fmt.Errorf("%w: unexpected status code (%d)", ErrAPIResponse, resp.StatusCode)
    }
}

// SetHTTPClient sets a custom HTTP client
func (c *Client) SetHTTPClient(client HTTPDoer) error {
	if client == nil {
		return fmt.Errorf("HTTP client cannot be nil")
	}
	c.httpClient = client
	return nil
}

package bento_test

import (
    "errors"
    "net/http"
    "testing"
    "time"

    bento "github.com/bentonow/bento-golang-sdk"
)

func TestClientCreation(t *testing.T) {
    tests := []struct {
        name        string
        config      *bento.Config
        expectError bool
        errorType   error
    }{
        {
            name: "valid config",
            config: &bento.Config{
                PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14b", // 32 chars
                SecretKey:      "s1803b8d410fd4ca3a7d1d1f5be6d3b65", // 32 chars
                SiteUUID:       "2103f23614d9877a6b4ee73d28a5c61d", // 32 chars
                Timeout:        10 * time.Second,
            },
            expectError: false,
        },
        {
            name: "missing publishable key",
            config: &bento.Config{
                SecretKey: "s1803b8d410fd4ca3a7d1d1f5be6d3b65",
                SiteUUID: "2103f23614d9877a6b4ee73d28a5c61d",
            },
            expectError: true,
            errorType:   bento.ErrInvalidConfig,
        },
        {
            name: "missing secret key",
            config: &bento.Config{
                PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14b",
                SiteUUID:       "2103f23614d9877a6b4ee73d28a5c61d",
            },
            expectError: true,
            errorType:   bento.ErrInvalidConfig,
        },
        {
            name: "missing site UUID",
            config: &bento.Config{
                PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14b",
                SecretKey:      "s1803b8d410fd4ca3a7d1d1f5be6d3b65",
            },
            expectError: true,
            errorType:   bento.ErrInvalidConfig,
        },
        {
            name: "with default timeout",
            config: &bento.Config{
                PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14b",
                SecretKey:      "s1803b8d410fd4ca3a7d1d1f5be6d3b65",
                SiteUUID:       "2103f23614d9877a6b4ee73d28a5c61d",
            },
            expectError: false,
        },
        {
            name: "invalid publishable key length",
            config: &bento.Config{
                PublishableKey: "tooshort",
                SecretKey:      "s1803b8d410fd4ca3a7d1d1f5be6d3b65",
                SiteUUID:       "2103f23614d9877a6b4ee73d28a5c61d",
            },
            expectError: true,
            errorType:   bento.ErrInvalidKeyLength,
        },
        {
            name: "invalid secret key length",
            config: &bento.Config{
                PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14b",
                SecretKey:      "tooshort",
                SiteUUID:       "2103f23614d9877a6b4ee73d28a5c61d",
            },
            expectError: true,
            errorType:   bento.ErrInvalidKeyLength,
        },
        {
            name: "invalid site UUID length",
            config: &bento.Config{
                PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14b",
                SecretKey:      "s1803b8d410fd4ca3a7d1d1f5be6d3b65",
                SiteUUID:       "tooshort",
            },
            expectError: true,
            errorType:   bento.ErrInvalidKeyLength,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, err := bento.NewClient(tt.config)

            if tt.expectError {
                if err == nil {
                    t.Error("expected error, got nil")
                    return
                }
                if tt.errorType != nil && !errors.Is(err, tt.errorType) {
                    t.Errorf("expected error type %v, got %v", tt.errorType, err)
                }
                return
            }
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
            if client == nil {
                t.Error("expected client, got nil")
            }
        })
    }
}

func TestClientHTTPClient(t *testing.T) {
    config := &bento.Config{
        PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14b",
        SecretKey:      "s1803b8d410fd4ca3a7d1d1f5be6d3b65",
        SiteUUID:       "2103f23614d9877a6b4ee73d28a5c61d",
        Timeout:        10 * time.Second,
    }

    client, err := bento.NewClient(config)
    if err != nil {
        t.Fatalf("failed to create client: %v", err)
    }

    // Test setting nil HTTP client
    err = client.SetHTTPClient(nil)
    if err == nil {
        t.Error("expected error when setting nil HTTP client, got nil")
    }

    // Test setting valid HTTP client
    mockClient := &mockHTTPClient{
        DoFunc: func(req *http.Request) (*http.Response, error) {
            return mockResponse(http.StatusOK, nil), nil
        },
    }
    err = client.SetHTTPClient(mockClient)
    if err != nil {
        t.Errorf("unexpected error setting valid HTTP client: %v", err)
    }
}
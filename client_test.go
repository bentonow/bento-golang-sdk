package bento_test

import (
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
	}{
		{
			name: "valid config",
			config: &bento.Config{
				PublishableKey: "test-key",
				SecretKey:      "test-secret",
				SiteUUID:       "test-uuid",
				Timeout:        10 * time.Second,
			},
			expectError: false,
		},
		{
			name: "missing publishable key",
			config: &bento.Config{
				SecretKey: "test-secret",
				SiteUUID:  "test-uuid",
			},
			expectError: true,
		},
		{
			name: "missing secret key",
			config: &bento.Config{
				PublishableKey: "test-key",
				SiteUUID:       "test-uuid",
			},
			expectError: true,
		},
		{
			name: "missing site UUID",
			config: &bento.Config{
				PublishableKey: "test-key",
				SecretKey:      "test-secret",
			},
			expectError: true,
		},
		{
			name: "with default timeout",
			config: &bento.Config{
				PublishableKey: "test-key",
				SecretKey:      "test-secret",
				SiteUUID:       "test-uuid",
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := bento.NewClient(tt.config)

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
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
	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		return mockResponse(http.StatusOK, map[string]string{"status": "ok"}), nil
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
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

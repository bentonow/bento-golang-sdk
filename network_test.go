package bento_test

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	bento "github.com/bentonow/bento-golang-sdk"
)

// TestNetworkFailures tests various network-related failures
func TestNetworkFailures(t *testing.T) {
	tests := []struct {
		name    string
		err     error
		handler func(*http.Request) (*http.Response, error)
	}{
		{
			name: "timeout error",
			err:  context.DeadlineExceeded,
			handler: func(req *http.Request) (*http.Response, error) {
				return nil, context.DeadlineExceeded
			},
		},
		{
			name: "network unreachable",
			err:  errors.New("network unreachable"),
			handler: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("network unreachable")
			},
		},
		{
			name: "connection reset",
			err:  errors.New("connection reset by peer"),
			handler: func(req *http.Request) (*http.Response, error) {
				return nil, errors.New("connection reset by peer")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(tt.handler)
			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			// Test various API endpoints
			_, err = client.GetTags(context.Background())
			if err == nil || err.Error() != tt.err.Error() {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}

			_, err = client.GetFields(context.Background())
			if err == nil || err.Error() != tt.err.Error() {
				t.Errorf("expected error %v, got %v", tt.err, err)
			}
		})
	}
}

// TestTimeoutConfiguration tests different timeout configurations
func TestTimeoutConfiguration(t *testing.T) {
	tests := []struct {
		name        string
		timeout     time.Duration
		expectError bool
	}{
		{
			name:        "zero timeout",
			timeout:     0,
			expectError: false, // Should use default
		},
		{
			name:        "negative timeout",
			timeout:     -1 * time.Second,
			expectError: true,
		},
		{
			name:        "very short timeout",
			timeout:     1 * time.Millisecond,
			expectError: false,
		},
		{
			name:        "very long timeout",
			timeout:     24 * time.Hour,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &bento.Config{
				PublishableKey: "test-key",
				SecretKey:      "test-secret",
				SiteUUID:       "test-uuid",
				Timeout:        tt.timeout,
			}

			_, err := bento.NewClient(config)
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestDataValidation tests various data validation scenarios
func TestDataValidation(t *testing.T) {
	tests := []struct {
		name        string
		input       interface{}
		expectError bool
	}{
		{
			name: "whitespace email",
			input: &bento.SubscriberInput{
				Email: "   ",
			},
			expectError: true,
		},
		{
			name: "malformed email",
			input: &bento.SubscriberInput{
				Email: "not-an-email",
			},
			expectError: true,
		},
		{
			name: "empty fields map",
			input: &bento.SubscriberInput{
				Email:  "test@example.com",
				Fields: map[string]interface{}{},
			},
			expectError: false,
		},
		{
			name: "nil fields map",
			input: &bento.SubscriberInput{
				Email:  "test@example.com",
				Fields: nil,
			},
			expectError: false,
		},
	}

	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		return mockResponse(http.StatusOK, map[string]string{"status": "success"}), nil
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var err error
			switch v := tt.input.(type) {
			case *bento.SubscriberInput:
				_, err = client.CreateSubscriber(context.Background(), v)
			}

			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

// TestContextBehaviors tests various context behaviors
func TestContextBehaviors(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func() (context.Context, context.CancelFunc)
		expectError bool
	}{
		{
			name: "immediate timeout",
			setupCtx: func() (context.Context, context.CancelFunc) {
				return context.WithTimeout(context.Background(), 0)
			},
			expectError: true,
		},
		{
			name: "parent cancellation",
			setupCtx: func() (context.Context, context.CancelFunc) {
				parentCtx, parentCancel := context.WithCancel(context.Background())
				parentCancel()
				return context.WithTimeout(parentCtx, time.Second)
			},
			expectError: true,
		},
		{
			name: "nested timeouts",
			setupCtx: func() (context.Context, context.CancelFunc) {
				ctx, _ := context.WithTimeout(context.Background(), time.Second)
				return context.WithTimeout(ctx, 500*time.Millisecond)
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, map[string]string{"status": "success"}), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			ctx, cancel := tt.setupCtx()
			defer cancel()

			_, err = client.GetTags(ctx)
			if tt.expectError && err == nil {
				t.Error("expected error, got nil")
			}
			if !tt.expectError && err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

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
    // Basic client configuration for testing
    config := &bento.Config{
        PublishableKey: "test-key",
        SecretKey:      "test-secret",
        SiteUUID:       "test-uuid",
        Timeout:        10 * time.Second,
    }

    tests := []struct {
        name    string
        err     error
        handler func(*http.Request) (*http.Response, error)
    }{
        {
            name: "timeout error",
            err:  context.DeadlineExceeded,
            handler: func(_ *http.Request) (*http.Response, error) {
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
            client, err := bento.NewClient(config)
            if err != nil {
                t.Fatalf("failed to create bento client: %v", err)
            }

            mockClient := &mockHTTPClient{
                DoFunc: tt.handler,
            }

            if err := client.SetHTTPClient(mockClient); err != nil {
                t.Fatalf("failed to set mock client: %v", err)
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

// TestContextBehaviors tests various context behaviors
func TestContextBehaviors(t *testing.T) {
    config := &bento.Config{
        PublishableKey: "test-key",
        SecretKey:      "test-secret",
        SiteUUID:       "test-uuid",
        Timeout:        10 * time.Second,
    }

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
                // Create parent context with longer timeout
                parentCtx, parentCancel := context.WithTimeout(context.Background(), 2*time.Second)
                // Create child context with shorter timeout
                childCtx, childCancel := context.WithTimeout(parentCtx, time.Second)

                // Return the child context and a combined cancel function
                return childCtx, func() {
                    childCancel()
                    parentCancel()
                }
            },
            expectError: false,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            client, err := bento.NewClient(config)
            if err != nil {
                t.Fatalf("failed to create bento client: %v", err)
            }

            mockClient := &mockHTTPClient{
                DoFunc: func(req *http.Request) (*http.Response, error) {
                    // Simulate a quick response for non-error cases
                    if !tt.expectError {
                        return mockResponse(http.StatusOK, map[string]string{"status": "success"}), nil
                    }

                    // Simulate delay for error cases
                    select {
                    case <-req.Context().Done():
                        return nil, req.Context().Err()
                    case <-time.After(100 * time.Millisecond):
                        return mockResponse(http.StatusOK, map[string]string{"status": "success"}), nil
                    }
                },
            }

            if err := client.SetHTTPClient(mockClient); err != nil {
                t.Fatalf("failed to set mock client: %v", err)
            }

            ctx, cancel := tt.setupCtx()
            defer cancel()

            _, err = client.GetTags(ctx)

            if tt.expectError {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                // Verify it's a context-related error
                if !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
                    t.Errorf("expected context error, got: %v", err)
                }
            } else {
                if err != nil {
                    t.Errorf("unexpected error: %v", err)
                }
            }
        })
    }
}
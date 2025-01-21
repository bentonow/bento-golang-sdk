package bento_test

import (
	"context"
	"net/http"
	"testing"
	"time"

	bento "github.com/bentonow/bento-golang-sdk"
)

func TestNestedTimeouts(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func() (context.Context, context.CancelFunc)
		expectError bool
	}{
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
		{
			name: "immediate parent timeout",
			setupCtx: func() (context.Context, context.CancelFunc) {
				parentCtx, parentCancel := context.WithTimeout(context.Background(), 0)
				childCtx, childCancel := context.WithTimeout(parentCtx, time.Second)
				return childCtx, func() {
					childCancel()
					parentCancel()
				}
			},
			expectError: true,
		},
		{
			name: "immediate child timeout",
			setupCtx: func() (context.Context, context.CancelFunc) {
				parentCtx, parentCancel := context.WithTimeout(context.Background(), 2*time.Second)
				childCtx, childCancel := context.WithTimeout(parentCtx, 0)
				return childCtx, func() {
					childCancel()
					parentCancel()
				}
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				select {
				case <-req.Context().Done():
					return nil, req.Context().Err()
				default:
					return mockResponse(http.StatusOK, map[string]interface{}{
						"data": []bento.TagData{},
					}), nil
				}
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			ctx, cancel := tt.setupCtx()
			defer cancel()

			// Use GetTags as a representative API call to test timeouts
			_, err = client.GetTags(ctx)

			if tt.expectError {
				if err == nil {
					t.Error("expected error due to timeout, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}
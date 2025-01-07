package bento_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	bento "github.com/bentonow/bento-golang-sdk"
)

func TestGetTags(t *testing.T) {
	// Sample tag data for responses
	sampleTags := []bento.TagData{
		{
			ID:   "tag1",
			Type: "tag",
			Attributes: struct {
				Name        string  `json:"name"`
				CreatedAt   string  `json:"created_at"`
				DiscardedAt *string `json:"discarded_at"`
				SiteID      int     `json:"site_id"`
			}{
				Name:      "test-tag-1",
				CreatedAt: time.Now().Format(time.RFC3339),
				SiteID:    1,
			},
		},
		{
			ID:   "tag2",
			Type: "tag",
			Attributes: struct {
				Name        string  `json:"name"`
				CreatedAt   string  `json:"created_at"`
				DiscardedAt *string `json:"discarded_at"`
				SiteID      int     `json:"site_id"`
			}{
				Name:      "test-tag-2",
				CreatedAt: time.Now().Format(time.RFC3339),
				SiteID:    1,
			},
		},
	}

	tests := []struct {
		name        string
		response    interface{}
		statusCode  int
		expectError bool
		wantTags    int
	}{
		{
			name: "successful retrieval",
			response: map[string]interface{}{
				"data": sampleTags,
			},
			statusCode:  http.StatusOK,
			expectError: false,
			wantTags:    2,
		},
		{
			name: "empty tags list",
			response: map[string]interface{}{
				"data": []bento.TagData{},
			},
			statusCode:  http.StatusOK,
			expectError: false,
			wantTags:    0,
		},
		{
			name:        "server error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
			wantTags:    0,
		},
		{
			name:        "invalid json response",
			response:    "{invalid-json",
			statusCode:  http.StatusOK,
			expectError: true,
			wantTags:    0,
		},
		{
			name: "unauthorized",
			response: map[string]interface{}{
				"error": "Unauthorized",
			},
			statusCode:  http.StatusUnauthorized,
			expectError: true,
			wantTags:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				// Verify authentication
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, map[string]string{
						"error": "Unauthorized",
					}), nil
				}

				// Verify request path and method
				if !strings.HasSuffix(req.URL.Path, "/fetch/tags") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", req.Method)
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			tags, err := client.GetTags(context.Background())

			// Verify error expectations
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify tags count
			if len(tags) != tt.wantTags {
				t.Errorf("got %d tags, want %d", len(tags), tt.wantTags)
			}
		})
	}
}

func TestCreateTag(t *testing.T) {
	sampleTag := bento.TagData{
		ID:   "new-tag-1",
		Type: "tag",
		Attributes: struct {
			Name        string  `json:"name"`
			CreatedAt   string  `json:"created_at"`
			DiscardedAt *string `json:"discarded_at"`
			SiteID      int     `json:"site_id"`
		}{
			Name:      "new-test-tag",
			CreatedAt: time.Now().Format(time.RFC3339),
			SiteID:    1,
		},
	}

	tests := []struct {
		name        string
		tagName     string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:    "successful creation",
			tagName: "new-test-tag",
			response: map[string]interface{}{
				"data": sampleTag,
			},
			statusCode:  http.StatusCreated,
			expectError: false,
		},
		{
			name:        "empty tag name",
			tagName:     "",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:    "server error",
			tagName: "test-tag",
			response: map[string]interface{}{
				"error": "Internal Server Error",
			},
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name:    "duplicate tag error",
			tagName: "existing-tag",
			response: map[string]interface{}{
				"error": "Tag already exists",
			},
			statusCode:  http.StatusConflict,
			expectError: true,
		},
		{
			name:        "invalid json response",
			tagName:     "test-tag",
			response:    "{invalid-json",
			statusCode:  http.StatusOK,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				// Verify authentication
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, map[string]string{
						"error": "Unauthorized",
					}), nil
				}

				// Verify request path and method
				if !strings.HasSuffix(req.URL.Path, "/fetch/tags") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", req.Method)
				}

				// Verify request body for non-empty tag names
				if tt.tagName != "" {
					body, err := io.ReadAll(req.Body)
					if err != nil {
						t.Fatalf("failed to read request body: %v", err)
					}

					var requestBody map[string]interface{}
					if err := json.Unmarshal(body, &requestBody); err != nil {
						t.Fatalf("invalid request body JSON: %v", err)
					}

					tag, ok := requestBody["tag"].(map[string]interface{})
					if !ok {
						t.Error("request body missing 'tag' object")
					}
					if tag["name"] != tt.tagName {
						t.Errorf("got tag name %v, want %v", tag["name"], tt.tagName)
					}
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			tag, err := client.CreateTag(context.Background(), tt.tagName)

			// Verify error expectations
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			// Verify tag data for successful creation
			if tag.Attributes.Name != tt.tagName {
				t.Errorf("got tag name %v, want %v", tag.Attributes.Name, tt.tagName)
			}
		})
	}
}

func TestGetTagsContextCancellation(t *testing.T) {
	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		// Simulate a delay to allow context cancellation to take effect
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

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	_, err = client.GetTags(ctx)
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

func TestCreateTagContextCancellation(t *testing.T) {
	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		// Simulate a delay to allow context cancellation to take effect
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
			return mockResponse(http.StatusCreated, map[string]interface{}{
				"data": bento.TagData{},
			}), nil
		}
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	_, err = client.CreateTag(ctx, "test-tag")
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

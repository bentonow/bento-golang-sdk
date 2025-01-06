package bento_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	bento "bento-golang-sdk"
)

func TestGetBroadcasts(t *testing.T) {
	sampleBroadcasts := []bento.BroadcastData{
		{
			Name:    "Test Broadcast 1",
			Subject: "Test Subject 1",
			Content: "<p>Test Content 1</p>",
			Type:    bento.BroadcastTypePlain,
			From: bento.ContactData{
				Name:  "Test Sender",
				Email: "sender@example.com",
			},
			BatchSizePerHour: 1000,
		},
		{
			Name:    "Test Broadcast 2",
			Subject: "Test Subject 2",
			Content: "<p>Test Content 2</p>",
			Type:    bento.BroadcastTypeRaw,
			From: bento.ContactData{
				Name:  "Test Sender",
				Email: "sender@example.com",
			},
			BatchSizePerHour: 2000,
		},
	}

	tests := []struct {
		name        string
		response    interface{}
		statusCode  int
		expectError bool
		wantCount   int
	}{
		{
			name: "successful retrieval",
			response: map[string]interface{}{
				"broadcasts": sampleBroadcasts,
			},
			statusCode:  http.StatusOK,
			expectError: false,
			wantCount:   2,
		},
		{
			name: "empty broadcasts list",
			response: map[string]interface{}{
				"broadcasts": []bento.BroadcastData{},
			},
			statusCode:  http.StatusOK,
			expectError: false,
			wantCount:   0,
		},
		{
			name:        "server error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
			wantCount:   0,
		},
		{
			name: "malformed response",
			response: map[string]interface{}{
				"broadcasts": make(chan int), // Will cause JSON marshal to fail
			},
			statusCode:  http.StatusOK,
			expectError: true,
			wantCount:   0,
		},
		{
			name:        "unauthorized",
			statusCode:  http.StatusUnauthorized,
			expectError: true,
			wantCount:   0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, map[string]string{
						"error": "Unauthorized",
					}), nil
				}

				if !strings.HasSuffix(req.URL.Path, "/fetch/broadcasts") {
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

			broadcasts, err := client.GetBroadcasts(context.Background())
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

			if len(broadcasts) != tt.wantCount {
				t.Errorf("got %d broadcasts, want %d", len(broadcasts), tt.wantCount)
			}
		})
	}
}

func TestCreateBroadcast(t *testing.T) {
	validBroadcasts := []bento.BroadcastData{
		{
			Name:    "Test Broadcast",
			Subject: "Test Subject",
			Content: "<p>Test Content</p>",
			Type:    bento.BroadcastTypePlain,
			From: bento.ContactData{
				Name:  "Test Sender",
				Email: "sender@example.com",
			},
			InclusiveTags:    "tag1,tag2",
			ExclusiveTags:    "tag3",
			SegmentID:        "segment123",
			BatchSizePerHour: 1000,
		},
	}

	tests := []struct {
		name        string
		broadcasts  []bento.BroadcastData
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:       "successful creation",
			broadcasts: validBroadcasts,
			response: map[string]interface{}{
				"status": "success",
			},
			statusCode:  http.StatusCreated,
			expectError: false,
		},
		{
			name:        "empty broadcasts list",
			broadcasts:  []bento.BroadcastData{},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "missing required fields",
			broadcasts: []bento.BroadcastData{
				{
					Name: "Test Broadcast",
					// Missing other required fields
				},
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "invalid email",
			broadcasts: []bento.BroadcastData{
				{
					Name:    "Test Broadcast",
					Subject: "Test Subject",
					Content: "<p>Test Content</p>",
					Type:    bento.BroadcastTypePlain,
					From: bento.ContactData{
						Name:  "Test Sender",
						Email: "invalid-email",
					},
					BatchSizePerHour: 1000,
				},
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "invalid batch size",
			broadcasts: []bento.BroadcastData{
				{
					Name:    "Test Broadcast",
					Subject: "Test Subject",
					Content: "<p>Test Content</p>",
					Type:    bento.BroadcastTypePlain,
					From: bento.ContactData{
						Name:  "Test Sender",
						Email: "sender@example.com",
					},
					BatchSizePerHour: 0, // Invalid batch size
				},
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:       "server error",
			broadcasts: validBroadcasts,
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name: "with all optional fields",
			broadcasts: []bento.BroadcastData{
				{
					Name:    "Test Broadcast",
					Subject: "Test Subject",
					Content: "<p>Test Content</p>",
					Type:    bento.BroadcastTypePlain,
					From: bento.ContactData{
						Name:  "Test Sender",
						Email: "sender@example.com",
					},
					InclusiveTags:    "tag1,tag2",
					ExclusiveTags:    "tag3",
					SegmentID:        "segment123",
					BatchSizePerHour: 1000,
				},
			},
			statusCode:  http.StatusCreated,
			expectError: false,
		},
		{
			name: "multiple broadcasts",
			broadcasts: []bento.BroadcastData{
				{
					Name:    "Test Broadcast 1",
					Subject: "Test Subject 1",
					Content: "<p>Test Content 1</p>",
					Type:    bento.BroadcastTypePlain,
					From: bento.ContactData{
						Email: "sender@example.com",
					},
					BatchSizePerHour: 1000,
				},
				{
					Name:    "Test Broadcast 2",
					Subject: "Test Subject 2",
					Content: "<p>Test Content 2</p>",
					Type:    bento.BroadcastTypeRaw,
					From: bento.ContactData{
						Email: "sender@example.com",
					},
					BatchSizePerHour: 2000,
				},
			},
			statusCode:  http.StatusCreated,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, map[string]string{
						"error": "Unauthorized",
					}), nil
				}

				if !strings.HasSuffix(req.URL.Path, "/batch/broadcasts") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", req.Method)
				}

				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}

				var requestBody map[string]interface{}
				if err := json.Unmarshal(body, &requestBody); err != nil {
					t.Fatalf("invalid request body JSON: %v", err)
				}

				if _, ok := requestBody["broadcasts"]; !ok {
					t.Error("request body missing 'broadcasts' field")
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			err = client.CreateBroadcast(context.Background(), tt.broadcasts)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestBroadcastContextCancellation(t *testing.T) {
	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
			return mockResponse(http.StatusOK, map[string]interface{}{
				"broadcasts": []bento.BroadcastData{},
			}), nil
		}
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	// Test GetBroadcasts with cancelled context
	_, err = client.GetBroadcasts(ctx)
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	// Test CreateBroadcast with cancelled context
	err = client.CreateBroadcast(ctx, []bento.BroadcastData{
		{
			Name:    "Test Broadcast",
			Subject: "Test Subject",
			Content: "<p>Test Content</p>",
			Type:    bento.BroadcastTypePlain,
			From: bento.ContactData{
				Email: "sender@example.com",
			},
			BatchSizePerHour: 1000,
		},
	})
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}
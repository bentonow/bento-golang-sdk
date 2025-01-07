package bento_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/bentonow/bento-golang-sdk"
)

func TestTrackEvent(t *testing.T) {
	validEvents := []bento.EventData{
		{
			Type:  "$completed_onboarding",
			Email: "test@example.com",
			Fields: map[string]interface{}{
				"first_name": "John",
				"last_name":  "Doe",
			},
		},
	}

	tests := []struct {
		name        string
		events      []bento.EventData
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:   "successful event tracking",
			events: validEvents,
			response: map[string]interface{}{
				"results": 1,
				"failed":  0,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:   "partial failure",
			events: validEvents,
			response: map[string]interface{}{
				"results": 0,
				"failed":  1,
			},
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name:        "empty events",
			events:      []bento.EventData{},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "invalid email",
			events: []bento.EventData{{
				Type:  "$completed_onboarding",
				Email: "invalid-email",
			}},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "missing event type",
			events: []bento.EventData{{
				Email: "test@example.com",
			}},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "with optional fields",
			events: []bento.EventData{{
				Type:  "$completed_onboarding",
				Email: "test@example.com",
				Fields: map[string]interface{}{
					"company": "Test Co",
					"role":    "Developer",
				},
				Details: map[string]interface{}{
					"source":  "API",
					"version": "1.0",
				},
			}},
			response: map[string]interface{}{
				"results": 1,
				"failed":  0,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name: "multiple events",
			events: []bento.EventData{
				{
					Type:  "$completed_onboarding",
					Email: "test1@example.com",
				},
				{
					Type:  "$viewed_docs",
					Email: "test2@example.com",
				},
			},
			response: map[string]interface{}{
				"results": 2,
				"failed":  0,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "server error",
			events:      validEvents,
			statusCode:  http.StatusInternalServerError,
			expectError: true,
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

				if !strings.HasSuffix(req.URL.Path, "/batch/events") {
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

				if _, ok := requestBody["events"]; !ok {
					t.Error("request body missing 'events' field")
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			err = client.TrackEvent(context.Background(), tt.events)
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

func TestTrackEventValidation(t *testing.T) {
	tests := []struct {
		name        string
		event       bento.EventData
		expectError bool
		errorMsg    string
	}{
		{
			name: "valid event",
			event: bento.EventData{
				Type:  "$test_event",
				Email: "test@example.com",
			},
			expectError: false,
		},
		{
			name: "invalid email",
			event: bento.EventData{
				Type:  "$test_event",
				Email: "not-an-email",
			},
			expectError: true,
			errorMsg:    "invalid email",
		},
		{
			name: "empty type",
			event: bento.EventData{
				Email: "test@example.com",
			},
			expectError: true,
			errorMsg:    "event type is required",
		},
		{
			name: "empty email",
			event: bento.EventData{
				Type: "$test_event",
			},
			expectError: true,
			errorMsg:    "invalid email",
		},
		{
			name: "with valid fields",
			event: bento.EventData{
				Type:  "$test_event",
				Email: "test@example.com",
				Fields: map[string]interface{}{
					"string_field": "value",
					"int_field":    42,
					"bool_field":   true,
				},
			},
			expectError: false,
		},
		{
			name: "with valid details",
			event: bento.EventData{
				Type:  "$test_event",
				Email: "test@example.com",
				Details: map[string]interface{}{
					"platform": "web",
					"version":  "1.0.0",
				},
			},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				return mockResponse(http.StatusOK, map[string]interface{}{
					"results": 1,
					"failed":  0,
				}), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			err = client.TrackEvent(context.Background(), []bento.EventData{tt.event})

			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				if err != nil && tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error containing '%s', got '%v'", tt.errorMsg, err)
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
	}
}

func TestTrackEventContextCancellation(t *testing.T) {
	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		// Simulate a delay to allow context cancellation to take effect
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
			return mockResponse(http.StatusOK, map[string]interface{}{
				"results": 1,
				"failed":  0,
			}), nil
		}
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	events := []bento.EventData{
		{
			Type:  "$test_event",
			Email: "test@example.com",
		},
	}

	err = client.TrackEvent(ctx, events)
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

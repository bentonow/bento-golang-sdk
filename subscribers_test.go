package bento_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"

	bento "bento-golang-sdk"
)

func TestFindSubscriber(t *testing.T) {
	tests := []struct {
		name        string
		email       string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:  "successful find",
			email: "test@example.com",
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "sub_123",
					"type": "subscriber",
					"attributes": map[string]interface{}{
						"uuid":  "uuid_123",
						"email": "test@example.com",
						"fields": map[string]interface{}{
							"first_name": "John",
						},
						"cached_tag_ids": []string{"tag1", "tag2"},
					},
				},
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "invalid email",
			email:       "invalid-email",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:  "subscriber not found",
			email: "notfound@example.com",
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"id": "",
				},
			},
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name:        "server error",
			email:       "test@example.com",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name:        "invalid response format",
			email:       "test@example.com",
			response:    "invalid json",
			statusCode:  http.StatusOK,
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

				if !strings.HasSuffix(req.URL.Path, "/fetch/subscribers") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", req.Method)
				}
				if req.URL.Query().Get("email") != tt.email {
					t.Errorf("unexpected email in query: %s", req.URL.Query().Get("email"))
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			subscriber, err := client.FindSubscriber(context.Background(), tt.email)
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
			if subscriber == nil {
				t.Error("expected subscriber, got nil")
				return
			}
			if subscriber.Attributes.Email != tt.email {
				t.Errorf("got email %s, want %s", subscriber.Attributes.Email, tt.email)
			}
		})
	}
}

func TestCreateSubscriber(t *testing.T) {
	validInput := &bento.SubscriberInput{
		Email:     "test@example.com",
		FirstName: "John",
		LastName:  "Doe",
		Tags:      []string{"tag1", "tag2"},
		Fields: map[string]interface{}{
			"company": "Test Co",
		},
	}

	tests := []struct {
		name        string
		input       *bento.SubscriberInput
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:  "successful creation",
			input: validInput,
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"id":   "sub_123",
					"type": "subscriber",
					"attributes": map[string]interface{}{
						"uuid":  "uuid_123",
						"email": "test@example.com",
					},
				},
			},
			statusCode:  http.StatusCreated,
			expectError: false,
		},
		{
			name: "invalid email",
			input: &bento.SubscriberInput{
				Email: "invalid-email",
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			input:       validInput,
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name:        "invalid response format",
			input:       validInput,
			response:    "invalid json",
			statusCode:  http.StatusOK,
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

				if !strings.HasSuffix(req.URL.Path, "/fetch/subscribers") {
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

				if _, ok := requestBody["subscriber"]; !ok {
					t.Error("request body missing 'subscriber' field")
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			subscriber, err := client.CreateSubscriber(context.Background(), tt.input)
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
			if subscriber == nil {
				t.Error("expected subscriber, got nil")
			}
		})
	}
}

func TestImportSubscribers(t *testing.T) {
	validSubscribers := []*bento.SubscriberInput{
		{
			Email:     "test1@example.com",
			FirstName: "John",
			LastName:  "Doe",
			Tags:      []string{"tag1"},
			Fields: map[string]interface{}{
				"company": "Test Co",
			},
		},
		{
			Email:     "test2@example.com",
			FirstName: "Jane",
			LastName:  "Doe",
			Tags:      []string{"tag2"},
		},
	}

	tests := []struct {
		name        string
		subscribers []*bento.SubscriberInput
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:        "successful import",
			subscribers: validSubscribers,
			response: map[string]interface{}{
				"results": 2,
				"failed":  0,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "empty subscribers list",
			subscribers: []*bento.SubscriberInput{},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "invalid email",
			subscribers: []*bento.SubscriberInput{
				{
					Email: "invalid-email",
				},
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			subscribers: validSubscribers,
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name:        "partial failure",
			subscribers: validSubscribers,
			response: map[string]interface{}{
				"results": 1,
				"failed":  1,
			},
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name:        "invalid response format",
			subscribers: validSubscribers,
			response:    "invalid json",
			statusCode:  http.StatusOK,
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

				if !strings.HasSuffix(req.URL.Path, "/batch/subscribers") {
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

				if _, ok := requestBody["subscribers"]; !ok {
					t.Error("request body missing 'subscribers' field")
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			err = client.ImportSubscribers(context.Background(), tt.subscribers)
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

func TestSubscriberWithContext(t *testing.T) {
	tests := []struct {
		name    string
		testFn  func(context.Context, *bento.Client) error
		ctxType string
	}{
		{
			name: "find subscriber with cancelled context",
			testFn: func(ctx context.Context, client *bento.Client) error {
				_, err := client.FindSubscriber(ctx, "test@example.com")
				return err
			},
			ctxType: "cancelled",
		},
		{
			name: "create subscriber with timeout context",
			testFn: func(ctx context.Context, client *bento.Client) error {
				_, err := client.CreateSubscriber(ctx, &bento.SubscriberInput{
					Email: "test@example.com",
				})
				return err
			},
			ctxType: "timeout",
		},
		{
			name: "import subscribers with deadline context",
			testFn: func(ctx context.Context, client *bento.Client) error {
				return client.ImportSubscribers(ctx, []*bento.SubscriberInput{
					{Email: "test@example.com"},
				})
			},
			ctxType: "deadline",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				// Simulate a delay to allow context cancellation to take effect
				select {
				case <-req.Context().Done():
					return nil, req.Context().Err()
				default:
					return mockResponse(http.StatusOK, map[string]interface{}{
						"results": 1,
					}), nil
				}
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			var ctx context.Context
			var cancel context.CancelFunc

			switch tt.ctxType {
			case "cancelled":
				ctx, cancel = context.WithCancel(context.Background())
				cancel()
			case "timeout":
				ctx, cancel = context.WithTimeout(context.Background(), 0)
				defer cancel()
			case "deadline":
				ctx, cancel = context.WithDeadline(context.Background(), time.Now())
				defer cancel()
			}

			err = tt.testFn(ctx, client)
			if err == nil {
				t.Error("expected error due to context, got nil")
			}
		})
	}
}
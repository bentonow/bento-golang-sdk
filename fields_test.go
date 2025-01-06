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

func TestGetFields(t *testing.T) {
	validResponse := bento.FieldsResponse{
		Data: []bento.FieldData{
			{
				ID:   "field_123",
				Type: "field",
				Attributes: bento.FieldAttributes{
					Name:        "Test Field",
					Key:         "test_field",
					Whitelisted: boolPtr(true),
					CreatedAt:   time.Now(),
				},
			},
		},
	}

	tests := []struct {
		name        string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:        "successful fields retrieval",
			response:    validResponse,
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "server error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name: "invalid response format",
			response: map[string]interface{}{
				"data": 123, // Invalid data type for Data field
			},
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name: "malformed json response",
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"invalid_field": make(chan int), // This will cause JSON marshal to fail
				},
			},
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name: "empty fields list",
			response: bento.FieldsResponse{
				Data: []bento.FieldData{},
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "unauthorized access",
			statusCode:  http.StatusUnauthorized,
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

				if !strings.HasSuffix(req.URL.Path, "/fetch/fields") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", req.Method)
				}

				if tt.name == "malformed json response" {
					// Create a custom response that will fail JSON marshaling
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader("{")), // Invalid JSON
						Header:     make(http.Header),
					}, nil
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			fields, err := client.GetFields(context.Background())
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

			if tt.response != nil {
				expectedFields, ok := tt.response.(bento.FieldsResponse)
				if ok && len(fields) != len(expectedFields.Data) {
					t.Errorf("got %d fields, want %d", len(fields), len(expectedFields.Data))
				}
			}
		})
	}
}

func TestCreateField(t *testing.T) {
	validResponse := struct {
		Data bento.FieldData `json:"data"`
	}{
		Data: bento.FieldData{
			ID:   "field_123",
			Type: "field",
			Attributes: bento.FieldAttributes{
				Name:        "Test Field",
				Key:         "test_field",
				Whitelisted: boolPtr(true),
				CreatedAt:   time.Now(),
			},
		},
	}

	tests := []struct {
		name        string
		key         string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:        "successful field creation",
			key:         "test_field",
			response:    validResponse,
			statusCode:  http.StatusCreated,
			expectError: false,
		},
		{
			name:        "empty field key",
			key:         "",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			key:         "test_field",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name: "invalid response format",
			key:  "test_field",
			response: map[string]interface{}{
				"data": 123, // Invalid data type for Data field
			},
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name: "malformed json response",
			key:  "test_field",
			response: map[string]interface{}{
				"data": map[string]interface{}{
					"invalid_field": make(chan int), // This will cause JSON marshal to fail
				},
			},
			statusCode:  http.StatusOK,
			expectError: true,
		},
		{
			name:        "unauthorized access",
			key:         "test_field",
			statusCode:  http.StatusUnauthorized,
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

				if !strings.HasSuffix(req.URL.Path, "/fetch/fields") {
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

				if field, ok := requestBody["field"].(map[string]interface{}); !ok {
					t.Error("request body missing 'field' object")
				} else if field["key"] != tt.key {
					t.Errorf("unexpected field key: got %v, want %v", field["key"], tt.key)
				}

				if tt.name == "malformed json response" {
					// Create a custom response that will fail JSON marshaling
					return &http.Response{
						StatusCode: tt.statusCode,
						Body:       io.NopCloser(strings.NewReader("{")), // Invalid JSON
						Header:     make(http.Header),
					}, nil
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			field, err := client.CreateField(context.Background(), tt.key)
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

			if field == nil {
				t.Error("expected field, got nil")
			}
		})
	}
}

func TestFieldsContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel context immediately

	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
			return mockResponse(http.StatusOK, bento.FieldsResponse{}), nil
		}
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	// Test GetFields with cancelled context
	_, err = client.GetFields(ctx)
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}

	// Test CreateField with cancelled context
	_, err = client.CreateField(ctx, "test_field")
	if err == nil {
		t.Error("expected error due to cancelled context, got nil")
	}
	if err != context.Canceled {
		t.Errorf("expected context.Canceled error, got %v", err)
	}
}

// Helper function to create a pointer to a bool
func boolPtr(b bool) *bool {
	return &b
}
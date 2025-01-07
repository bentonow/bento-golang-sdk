package bento_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	bento "github.com/bentonow/bento-golang-sdk"
)

func TestGetBlacklistStatus(t *testing.T) {
	tests := []struct {
		name        string
		data        *bento.BlacklistData
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name: "successful domain check",
			data: &bento.BlacklistData{
				Domain: "example.com",
			},
			response: map[string]interface{}{
				"status": "clean",
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name: "successful IP check",
			data: &bento.BlacklistData{
				IPAddress: "1.1.1.1",
			},
			response: map[string]interface{}{
				"status": "clean",
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name: "successful both checks",
			data: &bento.BlacklistData{
				Domain:    "example.com",
				IPAddress: "1.1.1.1",
			},
			response: map[string]interface{}{
				"status": "clean",
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "empty request",
			data:        &bento.BlacklistData{},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "invalid IP address",
			data: &bento.BlacklistData{
				IPAddress: "invalid-ip",
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			data:        &bento.BlacklistData{Domain: "example.com"},
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, nil), nil
				}

				if !strings.HasSuffix(req.URL.Path, "/experimental/blacklist.json") {
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

			result, err := client.GetBlacklistStatus(context.Background(), tt.data)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		name        string
		data        *bento.ValidationData
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name: "successful validation",
			data: &bento.ValidationData{
				EmailAddress: "test@example.com",
				FullName:     "John Doe",
				UserAgent:    "test-agent",
				IPAddress:    "1.1.1.1",
			},
			response: map[string]interface{}{
				"valid": true,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name: "successful validation without optional fields",
			data: &bento.ValidationData{
				EmailAddress: "test@example.com",
			},
			response: map[string]interface{}{
				"valid": true,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name: "invalid email format",
			data: &bento.ValidationData{
				EmailAddress: "invalid-email",
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "invalid IP address",
			data: &bento.ValidationData{
				EmailAddress: "test@example.com",
				IPAddress:    "invalid-ip",
			},
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name: "server error",
			data: &bento.ValidationData{
				EmailAddress: "test@example.com",
			},
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, nil), nil
				}

				if !strings.HasSuffix(req.URL.Path, "/experimental/validation") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", req.Method)
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			result, err := client.ValidateEmail(context.Background(), tt.data)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

func TestGetContentModeration(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:    "successful moderation",
			content: "This is a test message",
			response: map[string]interface{}{
				"status": "clean",
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "empty content",
			content:     "",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			content:     "Test content",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, nil), nil
				}

				if !strings.HasSuffix(req.URL.Path, "/experimental/content_moderation") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", req.Method)
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			result, err := client.GetContentModeration(context.Background(), tt.content)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

func TestGetGender(t *testing.T) {
	tests := []struct {
		name        string
		fullName    string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:     "successful gender prediction",
			fullName: "John Doe",
			response: map[string]interface{}{
				"gender":      "male",
				"probability": 0.95,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "empty name",
			fullName:    "",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			fullName:    "John Doe",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, nil), nil
				}

				if !strings.HasSuffix(req.URL.Path, "/experimental/gender") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", req.Method)
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			result, err := client.GetGender(context.Background(), tt.fullName)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

func TestGeoLocateIP(t *testing.T) {
	tests := []struct {
		name        string
		ipAddress   string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:      "successful geolocation",
			ipAddress: "1.1.1.1",
			response: map[string]interface{}{
				"country":   "US",
				"city":      "Los Angeles",
				"latitude":  34.0522,
				"longitude": -118.2437,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "invalid IP address",
			ipAddress:   "invalid-ip",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			ipAddress:   "1.1.1.1",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				if !validateAuthHeaders(req) {
					return mockResponse(http.StatusUnauthorized, nil), nil
				}

				if !strings.HasSuffix(req.URL.Path, "/experimental/geolocation") {
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

			result, err := client.GeoLocateIP(context.Background(), tt.ipAddress)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if result == nil {
				t.Error("expected result, got nil")
			}
		})
	}
}

func TestContextCancellation(t *testing.T) {
	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		select {
		case <-req.Context().Done():
			return nil, req.Context().Err()
		default:
			return mockResponse(http.StatusOK, map[string]interface{}{"status": "success"}), nil
		}
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	// Test each method with cancelled context
	t.Run("BlacklistStatus with cancelled context", func(t *testing.T) {
		_, err := client.GetBlacklistStatus(ctx, &bento.BlacklistData{Domain: "example.com"})
		if err == nil {
			t.Error("expected error due to cancelled context, got nil")
		}
	})

	t.Run("ValidateEmail with cancelled context", func(t *testing.T) {
		_, err := client.ValidateEmail(ctx, &bento.ValidationData{EmailAddress: "test@example.com"})
		if err == nil {
			t.Error("expected error due to cancelled context, got nil")
		}
	})

	t.Run("ContentModeration with cancelled context", func(t *testing.T) {
		_, err := client.GetContentModeration(ctx, "test content")
		if err == nil {
			t.Error("expected error due to cancelled context, got nil")
		}
	})

	t.Run("Gender prediction with cancelled context", func(t *testing.T) {
		_, err := client.GetGender(ctx, "John Doe")
		if err == nil {
			t.Error("expected error due to cancelled context, got nil")
		}
	})

	t.Run("GeoLocateIP with cancelled context", func(t *testing.T) {
		_, err := client.GeoLocateIP(ctx, "1.1.1.1")
		if err == nil {
			t.Error("expected error due to cancelled context, got nil")
		}
	})
}

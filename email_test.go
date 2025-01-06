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

func TestCreateEmails(t *testing.T) {
	validEmail := []bento.EmailData{
		{
			To:            "recipient@example.com",
			From:          "sender@example.com",
			Subject:       "Test Subject",
			HTMLBody:      "<p>Test Content</p>",
			Transactional: true,
			Personalizations: map[string]interface{}{
				"name": "John Doe",
			},
		},
	}

	tests := []struct {
		name        string
		emails      []bento.EmailData
		response    interface{}
		statusCode  int
		expectError bool
		wantResults int
	}{
		{
			name:   "successful email creation",
			emails: validEmail,
			response: map[string]interface{}{
				"results": 1,
			},
			statusCode:  http.StatusOK,
			expectError: false,
			wantResults: 1,
		},
		{
			name:        "empty emails list",
			emails:      []bento.EmailData{},
			statusCode:  http.StatusBadRequest,
			expectError: true,
			wantResults: 0,
		},
		{
			name: "invalid recipient email",
			emails: []bento.EmailData{{
				To:            "invalid-email",
				From:          "sender@example.com",
				Subject:       "Test Subject",
				HTMLBody:      "<p>Test Content</p>",
				Transactional: true,
			}},
			statusCode:  http.StatusBadRequest,
			expectError: true,
			wantResults: 0,
		},
		{
			name: "invalid sender email",
			emails: []bento.EmailData{{
				To:            "recipient@example.com",
				From:          "invalid-sender",
				Subject:       "Test Subject",
				HTMLBody:      "<p>Test Content</p>",
				Transactional: true,
			}},
			statusCode:  http.StatusBadRequest,
			expectError: true,
			wantResults: 0,
		},
		{
			name: "missing subject",
			emails: []bento.EmailData{{
				To:            "recipient@example.com",
				From:          "sender@example.com",
				HTMLBody:      "<p>Test Content</p>",
				Transactional: true,
			}},
			statusCode:  http.StatusBadRequest,
			expectError: true,
			wantResults: 0,
		},
		{
			name: "missing HTML body",
			emails: []bento.EmailData{{
				To:            "recipient@example.com",
				From:          "sender@example.com",
				Subject:       "Test Subject",
				Transactional: true,
			}},
			statusCode:  http.StatusBadRequest,
			expectError: true,
			wantResults: 0,
		},
		{
			name: "exceeds maximum batch size",
			emails: func() []bento.EmailData {
				emails := make([]bento.EmailData, 61) // Create 61 emails (exceeds 60 limit)
				for i := range emails {
					emails[i] = bento.EmailData{
						To:            "recipient@example.com",
						From:          "sender@example.com",
						Subject:       "Test Subject",
						HTMLBody:      "<p>Test Content</p>",
						Transactional: true,
					}
				}
				return emails
			}(),
			statusCode:  http.StatusBadRequest,
			expectError: true,
			wantResults: 0,
		},
		{
			name:   "server error",
			emails: validEmail,
			response: map[string]interface{}{
				"error": "Internal Server Error",
			},
			statusCode:  http.StatusInternalServerError,
			expectError: true,
			wantResults: 0,
		},
		{
			name: "with personalizations",
			emails: []bento.EmailData{{
				To:            "recipient@example.com",
				From:          "sender@example.com",
				Subject:       "Test Subject",
				HTMLBody:      "<p>Hello {{name}}!</p>",
				Transactional: true,
				Personalizations: map[string]interface{}{
					"name":    "John Doe",
					"company": "Acme Inc",
				},
			}},
			response: map[string]interface{}{
				"results": 1,
			},
			statusCode:  http.StatusOK,
			expectError: false,
			wantResults: 1,
		},
		{
			name: "multiple valid emails",
			emails: []bento.EmailData{
				{
					To:            "recipient1@example.com",
					From:          "sender@example.com",
					Subject:       "Test Subject 1",
					HTMLBody:      "<p>Test Content 1</p>",
					Transactional: true,
				},
				{
					To:            "recipient2@example.com",
					From:          "sender@example.com",
					Subject:       "Test Subject 2",
					HTMLBody:      "<p>Test Content 2</p>",
					Transactional: true,
				},
			},
			response: map[string]interface{}{
				"results": 2,
			},
			statusCode:  http.StatusOK,
			expectError: false,
			wantResults: 2,
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

				// Verify request path
				if !strings.HasSuffix(req.URL.Path, "/batch/emails") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}

				// Verify HTTP method
				if req.Method != http.MethodPost {
					t.Errorf("unexpected method: %s", req.Method)
				}

				// Verify request body format
				body, err := io.ReadAll(req.Body)
				if err != nil {
					t.Fatalf("failed to read request body: %v", err)
				}

				var requestBody map[string]interface{}
				if err := json.Unmarshal(body, &requestBody); err != nil {
					t.Fatalf("invalid request body JSON: %v", err)
				}

				if _, ok := requestBody["emails"]; !ok {
					t.Error("request body missing 'emails' field")
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			results, err := client.CreateEmails(context.Background(), tt.emails)

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

			// Verify results count
			if results != tt.wantResults {
				t.Errorf("got %d results, want %d", results, tt.wantResults)
			}
		})
	}
}

func TestCreateEmailsRequestValidation(t *testing.T) {
	client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
		return mockResponse(http.StatusOK, map[string]interface{}{
			"results": 1,
		}), nil
	})

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	// Test nil context
	_, err = client.CreateEmails(nil, []bento.EmailData{{
		To:       "test@example.com",
		From:     "sender@example.com",
		Subject:  "Test",
		HTMLBody: "<p>Test</p>",
	}})
	if err == nil {
		t.Error("expected error with nil context, got nil")
	}

	// Test with invalid personalizations
	_, err = client.CreateEmails(context.Background(), []bento.EmailData{{
		To:               "test@example.com",
		From:             "sender@example.com",
		Subject:          "Test",
		HTMLBody:         "<p>Test</p>",
		Personalizations: map[string]interface{}{
			"invalid": make(chan int), // Invalid type that can't be JSON marshaled
		},
	}})
	if err == nil {
		t.Error("expected error with invalid personalizations, got nil")
	}
}
package bento_test

import (
	"context"
	"net/http"
	"strings"
	"testing"

	bento "github.com/bentonow/bento-golang-sdk"
)

func TestGetSiteStats(t *testing.T) {
	tests := []struct {
		name        string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name: "successful stats retrieval",
			response: map[string]interface{}{
				"total_subscribers":  1000,
				"active_subscribers": 950,
				"growth_rate":        5.5,
				"engagement_rate":    75.2,
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "server error",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name: "malformed response",
			response: map[string]interface{}{
				"invalid": make(chan int), // This will cause JSON marshaling to fail
			},
			statusCode:  http.StatusOK,
			expectError: true,
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

				if !strings.HasSuffix(req.URL.Path, "/stats/site") {
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

			stats, err := client.GetSiteStats(context.Background())
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if stats == nil {
				t.Error("expected stats, got nil")
			}
		})
	}
}

func TestGetSegmentStats(t *testing.T) {
	tests := []struct {
		name        string
		segmentID   string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:      "successful segment stats retrieval",
			segmentID: "segment123",
			response: map[string]interface{}{
				"segment_size":    500,
				"conversion_rate": 25.5,
				"engagement_metrics": map[string]interface{}{
					"opens":  1500,
					"clicks": 750,
				},
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "empty segment ID",
			segmentID:   "",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			segmentID:   "segment123",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name:      "segment not found",
			segmentID: "nonexistent",
			response: map[string]interface{}{
				"error": "Segment not found",
			},
			statusCode:  http.StatusNotFound,
			expectError: true,
		},
		{
			name:      "malformed response",
			segmentID: "segment123",
			response: map[string]interface{}{
				"invalid": make(chan int), // This will cause JSON marshaling to fail
			},
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

				if !strings.HasSuffix(req.URL.Path, "/stats/segment") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", req.Method)
				}

				// Verify segment_id query parameter
				if req.URL.Query().Get("segment_id") != tt.segmentID {
					t.Errorf("unexpected segment_id: %s", req.URL.Query().Get("segment_id"))
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			stats, err := client.GetSegmentStats(context.Background(), tt.segmentID)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if stats == nil {
				t.Error("expected stats, got nil")
			}
		})
	}
}

func TestGetReportStats(t *testing.T) {
	tests := []struct {
		name        string
		reportID    string
		response    interface{}
		statusCode  int
		expectError bool
	}{
		{
			name:     "successful report stats retrieval",
			reportID: "report123",
			response: map[string]interface{}{
				"report_metrics": map[string]interface{}{
					"total_views":  1000,
					"unique_views": 750,
					"average_time": 120,
				},
				"time_series_data": []map[string]interface{}{
					{"date": "2024-01-01", "views": 100},
					{"date": "2024-01-02", "views": 150},
				},
			},
			statusCode:  http.StatusOK,
			expectError: false,
		},
		{
			name:        "empty report ID",
			reportID:    "",
			statusCode:  http.StatusBadRequest,
			expectError: true,
		},
		{
			name:        "server error",
			reportID:    "report123",
			statusCode:  http.StatusInternalServerError,
			expectError: true,
		},
		{
			name:     "report not found",
			reportID: "nonexistent",
			response: map[string]interface{}{
				"error": "Report not found",
			},
			statusCode:  http.StatusNotFound,
			expectError: true,
		},
		{
			name:     "malformed response",
			reportID: "report123",
			response: map[string]interface{}{
				"invalid": make(chan int), // This will cause JSON marshaling to fail
			},
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

				if !strings.HasSuffix(req.URL.Path, "/stats/report") {
					t.Errorf("unexpected path: %s", req.URL.Path)
				}
				if req.Method != http.MethodGet {
					t.Errorf("unexpected method: %s", req.Method)
				}

				// Verify report_id query parameter
				if req.URL.Query().Get("report_id") != tt.reportID {
					t.Errorf("unexpected report_id: %s", req.URL.Query().Get("report_id"))
				}

				return mockResponse(tt.statusCode, tt.response), nil
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			stats, err := client.GetReportStats(context.Background(), tt.reportID)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Errorf("unexpected error: %v", err)
			}
			if stats == nil {
				t.Error("expected stats, got nil")
			}
		})
	}
}

func TestStatsContextCancellation(t *testing.T) {
	testCases := []struct {
		name     string
		testFunc func(context.Context, *bento.Client) error
	}{
		{
			name: "GetSiteStats",
			testFunc: func(ctx context.Context, client *bento.Client) error {
				_, err := client.GetSiteStats(ctx)
				return err
			},
		},
		{
			name: "GetSegmentStats",
			testFunc: func(ctx context.Context, client *bento.Client) error {
				_, err := client.GetSegmentStats(ctx, "segment123")
				return err
			},
		},
		{
			name: "GetReportStats",
			testFunc: func(ctx context.Context, client *bento.Client) error {
				_, err := client.GetReportStats(ctx, "report123")
				return err
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			client, err := setupTestClient(func(req *http.Request) (*http.Response, error) {
				// Simulate a delay to allow context cancellation to take effect
				select {
				case <-req.Context().Done():
					return nil, req.Context().Err()
				default:
					return mockResponse(http.StatusOK, map[string]interface{}{
						"data": "test",
					}), nil
				}
			})

			if err != nil {
				t.Fatalf("failed to setup test client: %v", err)
			}

			ctx, cancel := context.WithCancel(context.Background())
			cancel() // Cancel context immediately

			err = tc.testFunc(ctx, client)
			if err == nil {
				t.Error("expected error due to cancelled context, got nil")
			}
			if err != context.Canceled {
				t.Errorf("expected context.Canceled error, got %v", err)
			}
		})
	}
}

func TestRequestValidation(t *testing.T) {
	client, err := setupTestClient(func(_ *http.Request) (*http.Response, error) {
        return mockResponse(http.StatusOK, map[string]interface{}{
            "data": "test",
        }), nil
    })

	if err != nil {
		t.Fatalf("failed to setup test client: %v", err)
	}

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately to create invalid context

	// Test cancelled context cases
	t.Run("GetSiteStats with cancelled context", func(t *testing.T) {
		_, err := client.GetSiteStats(ctx)
		if err == nil {
			t.Error("expected error with cancelled context, got nil")
		}
		if err != context.Canceled {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})

	t.Run("GetSegmentStats with cancelled context", func(t *testing.T) {
		_, err := client.GetSegmentStats(ctx, "segment123")
		if err == nil {
			t.Error("expected error with cancelled context, got nil")
		}
		if err != context.Canceled {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})

	t.Run("GetReportStats with cancelled context", func(t *testing.T) {
		_, err := client.GetReportStats(ctx, "report123")
		if err == nil {
			t.Error("expected error with cancelled context, got nil")
		}
		if err != context.Canceled {
			t.Errorf("expected context.Canceled error, got %v", err)
		}
	})
}
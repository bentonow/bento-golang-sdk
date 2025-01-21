package bento_test

import (
    "encoding/json"
    "io"
    "net/http"
    "strings"
    "time"

    bento "github.com/bentonow/bento-golang-sdk"
)

// mockHTTPClient is a test helper that returns a custom http.Client
type mockHTTPClient struct {
    DoFunc func(req *http.Request) (*http.Response, error)
}

func (m *mockHTTPClient) Do(req *http.Request) (*http.Response, error) {
    return m.DoFunc(req)
}

// setupTestClient creates a new Client with mocked HTTP responses
func setupTestClient(handler func(req *http.Request) (*http.Response, error)) (*bento.Client, error) {
    config := &bento.Config{
        PublishableKey: "pc422f7e69255a4bf9c9fafcaac64b14", // 32 chars exactly
        SecretKey:      "s1803b8d410fd4ca3a7d1d1f5be6d3b6", // 32 chars exactly
        SiteUUID:       "2103f23614d9877a6b4ee73d28a5c610", // 32 chars exactly
        Timeout:        10 * time.Second,
    }

    client, err := bento.NewClient(config)
    if err != nil {
        return nil, err
    }

    if err := client.SetHTTPClient(&mockHTTPClient{DoFunc: handler}); err != nil {
        return nil, err
    }

    return client, nil
}

// mockResponse creates a mock HTTP response with the given status code and body
func mockResponse(statusCode int, body interface{}) *http.Response {
    jsonBody, _ := json.Marshal(body)
    return &http.Response{
        StatusCode: statusCode,
        Body:       io.NopCloser(strings.NewReader(string(jsonBody))),
        Header:     make(http.Header),
    }
}

// validateAuthHeaders checks if the request has proper authentication headers
func validateAuthHeaders(req *http.Request) bool {
    auth := req.Header.Get("Authorization")
    return auth != "" && strings.HasPrefix(auth, "Basic ")
}
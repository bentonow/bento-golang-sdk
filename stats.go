package bento

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// GetSiteStats retrieves site statistics
func (c *Client) GetSiteStats(ctx context.Context) (map[string]interface{}, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/stats/site", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetSegmentStats retrieves segment statistics
func (c *Client) GetSegmentStats(ctx context.Context, segmentID string) (map[string]interface{}, error) {
	if segmentID == "" {
		return nil, fmt.Errorf("%w: segment ID is required", ErrInvalidSegmentID)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/stats/segment", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("segment_id", segmentID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetReportStats retrieves report statistics
func (c *Client) GetReportStats(ctx context.Context, reportID string) (map[string]interface{}, error) {
	if reportID == "" {
		return nil, fmt.Errorf("%w: report ID is required", ErrInvalidRequest)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet,
		fmt.Sprintf("%s/stats/report", c.baseURL), nil)
	if err != nil {
		return nil, err
	}

	q := req.URL.Query()
	q.Add("report_id", reportID)
	req.URL.RawQuery = q.Encode()

	resp, err := c.do(req)
	if err != nil {
		return nil, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	return result, nil
}

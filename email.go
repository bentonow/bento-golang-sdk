package bento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

// CreateEmails sends one or more emails through Bento
func (c *Client) CreateEmails(ctx context.Context, emails []EmailData) (int, error) {
	if len(emails) == 0 {
		return 0, fmt.Errorf("%w: no emails provided", ErrInvalidRequest)
	}

	if len(emails) > 60 {
		return 0, fmt.Errorf("%w: maximum of 60 emails allowed per request", ErrInvalidRequest)
	}

	// Validate all emails before sending
	for _, email := range emails {
		if _, err := mail.ParseAddress(email.To); err != nil {
			return 0, fmt.Errorf("%w: invalid recipient email: %s", ErrInvalidEmail, email.To)
		}
		if _, err := mail.ParseAddress(email.From); err != nil {
			return 0, fmt.Errorf("%w: invalid sender email: %s", ErrInvalidEmail, email.From)
		}
		if email.Subject == "" {
			return 0, fmt.Errorf("%w: subject is required", ErrInvalidRequest)
		}
		if email.HTMLBody == "" {
			return 0, fmt.Errorf("%w: html_body is required", ErrInvalidRequest)
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"emails": emails,
	})
	if err != nil {
		return 0, err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/batch/emails", c.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return 0, err
	}

	resp, err := c.do(req)
	if err != nil {
		return 0, err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return 0, fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result struct {
		Results int `json:"results"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return 0, err
	}

	return result.Results, nil
}

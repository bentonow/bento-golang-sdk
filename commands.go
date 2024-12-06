package bento

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/mail"
)

// SubscriberCommand executes a command on a subscriber
func (c *Client) SubscriberCommand(ctx context.Context, commands []CommandData) error {
	if len(commands) == 0 {
		return ErrInvalidRequest
	}

	// Validate all commands before sending
	for _, cmd := range commands {
		if _, err := mail.ParseAddress(cmd.Email); err != nil {
			return fmt.Errorf("%w: %s", ErrInvalidEmail, cmd.Email)
		}
		if cmd.Query == "" {
			return fmt.Errorf("%w: command query is required", ErrInvalidRequest)
		}
		if err := validateCommandType(cmd.Command); err != nil {
			return err
		}
	}

	body, err := json.Marshal(map[string]interface{}{
		"command": commands,
	})
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost,
		fmt.Sprintf("%s/fetch/commands", c.baseURL), bytes.NewBuffer(body))
	if err != nil {
		return err
	}

	resp, err := c.do(req)
	if err != nil {
		return err
	}
	defer func() { _ = resp.Body.Close() }()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("%w: %d", ErrAPIResponse, resp.StatusCode)
	}

	var result struct {
		Results int `json:"results"`
		Failed  int `json:"failed"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return err
	}

	if result.Failed > 0 {
		return fmt.Errorf("command execution partially failed: %d succeeded, %d failed",
			result.Results, result.Failed)
	}

	return nil
}

// validateCommandType ensures the command type is valid
func validateCommandType(cmd CommandType) error {
	valid := map[CommandType]bool{
		CommandAddTag:         true,
		CommandAddTagViaEvent: true,
		CommandRemoveTag:      true,
		CommandAddField:       true,
		CommandRemoveField:    true,
		CommandSubscribe:      true,
		CommandUnsubscribe:    true,
		CommandChangeEmail:    true,
	}

	if !valid[cmd] {
		return fmt.Errorf("%w: invalid command type: %s", ErrInvalidRequest, cmd)
	}
	return nil
}

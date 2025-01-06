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

func TestSubscriberCommand(t *testing.T) {
    validCommands := []bento.CommandData{
        {
            Command: bento.CommandAddTag,
            Email:   "test@example.com",
            Query:   "new-tag",
        },
    }

    tests := []struct {
        name        string
        commands    []bento.CommandData
        response    interface{}
        statusCode  int
        expectError bool
    }{
        {
            name:     "successful command execution",
            commands: validCommands,
            response: map[string]interface{}{
                "results": 1,
                "failed":  0,
            },
            statusCode:  http.StatusOK,
            expectError: false,
        },
        {
            name:     "partial failure",
            commands: validCommands,
            response: map[string]interface{}{
                "results": 0,
                "failed":  1,
            },
            statusCode:  http.StatusOK,
            expectError: true,
        },
        {
            name:        "empty commands",
            commands:    []bento.CommandData{},
            statusCode:  http.StatusBadRequest,
            expectError: true,
        },
        {
            name: "invalid email",
            commands: []bento.CommandData{{
                Command: bento.CommandAddTag,
                Email:   "invalid-email",
                Query:   "test-tag",
            }},
            statusCode:  http.StatusBadRequest,
            expectError: true,
        },
        {
            name: "empty query",
            commands: []bento.CommandData{{
                Command: bento.CommandAddTag,
                Email:   "test@example.com",
                Query:   "",
            }},
            statusCode:  http.StatusBadRequest,
            expectError: true,
        },
        {
            name: "invalid command type",
            commands: []bento.CommandData{{
                Command: "invalid_command",
                Email:   "test@example.com",
                Query:   "test-tag",
            }},
            statusCode:  http.StatusBadRequest,
            expectError: true,
        },
        {
            name:        "server error",
            commands:    validCommands,
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

                if !strings.HasSuffix(req.URL.Path, "/fetch/commands") {
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

                if _, ok := requestBody["command"]; !ok {
                    t.Error("request body missing 'command' field")
                }

                return mockResponse(tt.statusCode, tt.response), nil
            })

            if err != nil {
                t.Fatalf("failed to setup test client: %v", err)
            }

            err = client.SubscriberCommand(context.Background(), tt.commands)
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

func TestValidateCommandType(t *testing.T) {
    tests := []struct {
        name        string
        commandType bento.CommandType
        expectError bool
    }{
        {
            name:        "valid add tag command",
            commandType: bento.CommandAddTag,
            expectError: false,
        },
        {
            name:        "valid add tag via event command",
            commandType: bento.CommandAddTagViaEvent,
            expectError: false,
        },
        {
            name:        "valid remove tag command",
            commandType: bento.CommandRemoveTag,
            expectError: false,
        },
        {
            name:        "valid add field command",
            commandType: bento.CommandAddField,
            expectError: false,
        },
        {
            name:        "valid remove field command",
            commandType: bento.CommandRemoveField,
            expectError: false,
        },
        {
            name:        "valid subscribe command",
            commandType: bento.CommandSubscribe,
            expectError: false,
        },
        {
            name:        "valid unsubscribe command",
            commandType: bento.CommandUnsubscribe,
            expectError: false,
        },
        {
            name:        "valid change email command",
            commandType: bento.CommandChangeEmail,
            expectError: false,
        },
        {
            name:        "invalid command type",
            commandType: "invalid_command",
            expectError: true,
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

            cmd := bento.CommandData{
                Command: tt.commandType,
                Email:   "test@example.com",
                Query:   "test-query",
            }

            err = client.SubscriberCommand(context.Background(), []bento.CommandData{cmd})

            if tt.expectError {
                if err == nil {
                    t.Error("expected error, got nil")
                }
                if err != nil && !strings.Contains(err.Error(), "invalid command type") {
                    t.Errorf("unexpected error message: %v", err)
                }
                return
            }
            if err != nil {
                t.Errorf("unexpected error: %v", err)
            }
        })
    }
}
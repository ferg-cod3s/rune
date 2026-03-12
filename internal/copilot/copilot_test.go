package copilot

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestCopilotRequest_Marshal(t *testing.T) {
	req := CopilotRequest{
		Messages: []map[string]string{
			{"role": "user", "content": "hello"},
		},
	}

	data, err := json.Marshal(req)
	if err != nil {
		t.Fatalf("failed to marshal CopilotRequest: %v", err)
	}

	var unmarshaled CopilotRequest
	if err := json.Unmarshal(data, &unmarshaled); err != nil {
		t.Fatalf("failed to unmarshal CopilotRequest: %v", err)
	}

	if len(unmarshaled.Messages) != 1 {
		t.Errorf("expected 1 message, got %d", len(unmarshaled.Messages))
	}
	if unmarshaled.Messages[0]["role"] != "user" {
		t.Errorf("expected role 'user', got %q", unmarshaled.Messages[0]["role"])
	}
	if unmarshaled.Messages[0]["content"] != "hello" {
		t.Errorf("expected content 'hello', got %q", unmarshaled.Messages[0]["content"])
	}
}

func TestCallCopilotLLM_Success(t *testing.T) {
	// Create a mock server that returns a successful response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify request method, path, and headers
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/chat/completions" {
			t.Errorf("expected path '/chat/completions', got %q", r.URL.Path)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Errorf("expected Content-Type application/json, got %s", r.Header.Get("Content-Type"))
		}
		if r.Header.Get("Authorization") != "Bearer test-token" {
			t.Errorf("expected Authorization 'Bearer test-token', got %s", r.Header.Get("Authorization"))
		}

		// Decode the request body
		var req CopilotRequest
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			t.Errorf("failed to decode request body: %v", err)
		}
		if len(req.Messages) != 1 || req.Messages[0]["content"] != "hello" {
			t.Errorf("unexpected request body: %+v", req)
		}

		// Return a successful response
		resp := CopilotResponse{
			Choices: []struct {
				Message struct {
					Content string `json:"content"`
				} `json:"message"`
			}{
				{Message: struct {
					Content string `json:"content"`
				}{Content: "Hello from Copilot!"}},
			},
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Override package-level variables to use the mock server
	origBase := copilotBaseURL
	origClient := copilotHTTPClient
	copilotBaseURL = server.URL
	copilotHTTPClient = server.Client()
	defer func() {
		copilotBaseURL = origBase
		copilotHTTPClient = origClient
	}()

	result, err := CallCopilotLLM("hello", "test-token")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result != "Hello from Copilot!" {
		t.Errorf("expected 'Hello from Copilot!', got %q", result)
	}
}

func TestCallCopilotLLM_EmptyChoices(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := CopilotResponse{Choices: nil}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			t.Errorf("failed to encode response: %v", err)
		}
	}))
	defer server.Close()

	// Override package-level variables to use the mock server
	origBase := copilotBaseURL
	origClient := copilotHTTPClient
	copilotBaseURL = server.URL
	copilotHTTPClient = server.Client()
	defer func() {
		copilotBaseURL = origBase
		copilotHTTPClient = origClient
	}()

	_, err := CallCopilotLLM("hello", "test-token")
	if err == nil {
		t.Fatal("expected error for empty choices, got nil")
	}
	if !errors.Is(err, ErrNoChoices) {
		t.Errorf("expected ErrNoChoices, got %v", err)
	}
}

func TestCopilotResponse_Unmarshal(t *testing.T) {
	jsonData := `{
		"choices": [
			{
				"message": {
					"content": "Test response content"
				}
			}
		]
	}`

	var resp CopilotResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("failed to unmarshal CopilotResponse: %v", err)
	}

	if len(resp.Choices) != 1 {
		t.Fatalf("expected 1 choice, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "Test response content" {
		t.Errorf("expected 'Test response content', got %q", resp.Choices[0].Message.Content)
	}
}

func TestCopilotResponse_MultipleChoices(t *testing.T) {
	jsonData := `{
		"choices": [
			{"message": {"content": "First choice"}},
			{"message": {"content": "Second choice"}}
		]
	}`

	var resp CopilotResponse
	if err := json.Unmarshal([]byte(jsonData), &resp); err != nil {
		t.Fatalf("failed to unmarshal CopilotResponse: %v", err)
	}

	if len(resp.Choices) != 2 {
		t.Fatalf("expected 2 choices, got %d", len(resp.Choices))
	}

	if resp.Choices[0].Message.Content != "First choice" {
		t.Errorf("expected 'First choice', got %q", resp.Choices[0].Message.Content)
	}
	if resp.Choices[1].Message.Content != "Second choice" {
		t.Errorf("expected 'Second choice', got %q", resp.Choices[1].Message.Content)
	}
}

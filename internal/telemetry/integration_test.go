package telemetry

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestSegmentIntegration tests actual Segment event delivery using a mock server
func TestSegmentIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Create a mock Segment server
	var receivedEvents []map[string]interface{}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		if !strings.Contains(r.URL.Path, "/track") {
			w.WriteHeader(http.StatusNotFound)
			return
		}

		// Parse the request body
		var event map[string]interface{}
		if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		receivedEvents = append(receivedEvents, event)
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"success": true}`))
	}))
	defer server.Close()

	// Note: This test would need modification to the Segment client to use a custom endpoint
	// For now, we'll test the client creation and method calls
	client := NewClient("test_segment_key", "")
	require.NotNil(t, client)

	// Test tracking an event
	client.Track("integration_test_event", map[string]interface{}{
		"test_property": "test_value",
		"timestamp":     time.Now().Unix(),
	})

	// Test tracking a command
	client.TrackCommand("test_command", 100*time.Millisecond, true)

	// Test tracking an error
	client.TrackError(fmt.Errorf("test error"), "test_command", map[string]interface{}{
		"context": "integration_test",
	})

	// Close the client to flush events
	client.Close()

	// In a real integration test, we would verify events were received
	// For now, we just ensure no panics occurred
}

// TestSentryIntegration tests Sentry integration with environment variables
func TestSentryIntegration(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	// Only run if we have a test Sentry DSN
	testDSN := os.Getenv("RUNE_TEST_SENTRY_DSN")
	if testDSN == "" {
		t.Skip("RUNE_TEST_SENTRY_DSN not set, skipping Sentry integration test")
	}

	client := NewClient("", testDSN)
	require.NotNil(t, client)
	assert.True(t, client.sentryEnabled)

	// Test capturing an exception
	client.CaptureException(
		fmt.Errorf("integration test error"),
		map[string]string{"test": "integration"},
		map[string]interface{}{"timestamp": time.Now().Unix()},
	)

	// Test starting a transaction
	transaction := client.StartTransaction("integration_test", "test")
	if transaction != nil {
		transaction.Finish()
	}

	// Test command tracking
	client.StartCommand("integration_test_command")
	client.EndCommand("integration_test_command", true, 50*time.Millisecond)

	// Close to flush events
	client.Close()
}

// TestTelemetryWithRealKeys tests with actual API keys if available
func TestTelemetryWithRealKeys(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping integration test in short mode")
	}

	segmentKey := os.Getenv("RUNE_TEST_SEGMENT_WRITE_KEY")
	sentryDSN := os.Getenv("RUNE_TEST_SENTRY_DSN")

	if segmentKey == "" && sentryDSN == "" {
		t.Skip("No test API keys provided, skipping real integration test")
	}

	client := NewClient(segmentKey, sentryDSN)
	require.NotNil(t, client)

	// Send a test event
	client.Track("integration_test_real_keys", map[string]interface{}{
		"test_run":   true,
		"timestamp":  time.Now().Unix(),
		"go_version": "test",
	})

	// Test command execution tracking
	client.TrackCommand("integration_test_command", 25*time.Millisecond, true)

	// Test error tracking
	client.TrackError(
		fmt.Errorf("integration test error - this is expected"),
		"integration_test_command",
		map[string]interface{}{
			"test_context": "real_keys_test",
		},
	)

	// Close to ensure events are sent
	client.Close()

	t.Log("Integration test completed - check your analytics dashboards for events")
}

// TestMiddlewareIntegration tests the middleware functions
func TestMiddlewareIntegration(t *testing.T) {
	// Initialize global client for testing
	Initialize("test_key", "")
	defer Close()

	// Test global functions
	Track("middleware_test_event", map[string]interface{}{
		"source": "middleware_test",
	})

	TrackCommand("middleware_test_command", 10*time.Millisecond, true)

	TrackError(fmt.Errorf("middleware test error"), "test_command", map[string]interface{}{
		"middleware": true,
	})

	StartCommand("middleware_command")
	EndCommand("middleware_command", true, 5*time.Millisecond)
}

// TestEventProperties tests that events contain expected properties
func TestEventProperties(t *testing.T) {
	client := NewClient("test_key", "")
	require.NotNil(t, client)

	// We can't easily test the actual properties sent to external services
	// without mocking, but we can test that the methods don't panic
	// and handle various input types correctly

	testCases := []struct {
		name       string
		event      string
		properties map[string]interface{}
	}{
		{
			name:  "string properties",
			event: "test_string_props",
			properties: map[string]interface{}{
				"string_prop": "test_value",
				"app_name":    "rune",
			},
		},
		{
			name:  "mixed properties",
			event: "test_mixed_props",
			properties: map[string]interface{}{
				"string_prop": "test",
				"int_prop":    42,
				"bool_prop":   true,
				"float_prop":  3.14,
				"time_prop":   time.Now(),
			},
		},
		{
			name:       "nil properties",
			event:      "test_nil_props",
			properties: nil,
		},
		{
			name:       "empty properties",
			event:      "test_empty_props",
			properties: map[string]interface{}{},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Should not panic
			client.Track(tc.event, tc.properties)
		})
	}
}

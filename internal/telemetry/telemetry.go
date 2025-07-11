package telemetry

import (
	"context"
	"fmt"
	"os"
	"runtime"
	"time"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/getsentry/sentry-go"
	"github.com/segmentio/analytics-go/v3"
)

type Client struct {
	segmentClient analytics.Client
	sentryEnabled bool
	enabled       bool
	userID        string
	sessionID     string
}

func NewClient(userSegmentKey, userSentryDSN string) *Client {
	// Check if telemetry is disabled via environment variable or config
	enabled := os.Getenv("RUNE_TELEMETRY_DISABLED") != "true"

	// Debug logging
	if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Telemetry enabled: %v\n", enabled)
		fmt.Printf("DEBUG: User Segment Key: %s\n", userSegmentKey)
		fmt.Printf("DEBUG: User Sentry DSN: %s\n", userSentryDSN)
		fmt.Printf("DEBUG: Build-time Segment Key: %s\n", segmentWriteKey)
		fmt.Printf("DEBUG: Build-time Sentry DSN: %s\n", sentryDSN)
	}

	// Generate or load user ID (anonymous)
	userID := getUserID()

	// Use build-time keys if user hasn't provided their own
	finalSegmentKey := userSegmentKey
	if finalSegmentKey == "" {
		finalSegmentKey = segmentWriteKey
	}

	finalSentryDSN := userSentryDSN
	if finalSentryDSN == "" {
		finalSentryDSN = sentryDSN
	}

	// Debug logging for final values
	if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Final Segment Key: %s\n", finalSegmentKey)
		fmt.Printf("DEBUG: Final Sentry DSN: %s\n", finalSentryDSN)
		fmt.Printf("DEBUG: User ID: %s\n", userID)
	}

	client := &Client{
		enabled:       enabled,
		userID:        userID,
		sentryEnabled: finalSentryDSN != "",
		sessionID:     generateSessionID(),
	}

	// Initialize Segment client if enabled and write key provided
	if enabled && finalSegmentKey != "" {
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Initializing Segment client with key: %s\n", finalSegmentKey)
		}
		segmentClient := analytics.New(finalSegmentKey)
		client.segmentClient = segmentClient
	} else if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Segment client not initialized - enabled: %v, key: %s\n", enabled, finalSegmentKey)
	}

	// Initialize Sentry if enabled and DSN provided
	if enabled && finalSentryDSN != "" {
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Initializing Sentry with DSN: %s\n", finalSentryDSN)
		}
		err := sentry.Init(sentry.ClientOptions{
			Dsn:              finalSentryDSN,
			Environment:      getEnvironment(),
			Release:          getVersion(),
			AttachStacktrace: true,
			BeforeSend: func(event *sentry.Event, hint *sentry.EventHint) *sentry.Event {
				// Add user context
				event.User = sentry.User{
					ID: userID,
				}
				// Add app context
				event.Contexts["app"] = map[string]interface{}{
					"name":    "rune",
					"version": getVersion(),
				}
				event.Contexts["os"] = map[string]interface{}{
					"name":    runtime.GOOS,
					"version": getOSVersion(),
				}
				return event
			},
		})
		if err != nil {
			if os.Getenv("RUNE_DEBUG") == "true" {
				fmt.Printf("DEBUG: Sentry initialization failed: %v\n", err)
			}
			// Silently fail for telemetry initialization
			client.sentryEnabled = false
		} else if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Sentry initialized successfully\n")
		}
	} else if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Sentry not initialized - enabled: %v, DSN: %s\n", enabled, finalSentryDSN)
	}

	return client
}

func (c *Client) Track(event string, properties map[string]interface{}) {
	if !c.enabled {
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Telemetry disabled, not tracking event: %s\n", event)
		}
		return
	}

	if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Tracking event: %s\n", event)
	}

	// Add default properties
	if properties == nil {
		properties = make(map[string]interface{})
	}

	// Add system context
	properties["app_name"] = "rune"
	properties["app_version"] = getVersion()
	properties["os_name"] = runtime.GOOS
	properties["os_version"] = getOSVersion()

	// Send to Segment if available
	if c.segmentClient != nil {
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Sending to Segment: %s\n", event)
		}
		err := c.segmentClient.Enqueue(analytics.Track{
			UserId:     c.userID,
			Event:      event,
			Properties: properties,
			Timestamp:  time.Now(),
		})
		if err != nil && os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Segment enqueue error: %v\n", err)
		}
	} else if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Segment client not available for event: %s\n", event)
	}

	// Send to Sentry as breadcrumb for context
	if c.sentryEnabled {
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Adding Sentry breadcrumb: %s\n", event)
		}
		sentry.AddBreadcrumb(&sentry.Breadcrumb{
			Message:   event,
			Category:  "telemetry",
			Level:     sentry.LevelInfo,
			Data:      properties,
			Timestamp: time.Now(),
		})
	} else if os.Getenv("RUNE_DEBUG") == "true" {
		fmt.Printf("DEBUG: Sentry not enabled for event: %s\n", event)
	}
}

func (c *Client) TrackError(err error, command string, properties map[string]interface{}) {
	if !c.enabled {
		return
	}

	if properties == nil {
		properties = make(map[string]interface{})
	}

	properties["error"] = err.Error()
	properties["command"] = command
	properties["error_type"] = fmt.Sprintf("%T", err)

	// Track error event in Segment
	c.Track("error", properties)

	// Send to Sentry for error tracking
	if c.sentryEnabled {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("command", command)
			scope.SetContext("error_details", properties)
			sentry.CaptureException(err)
		})
	}
}

func (c *Client) TrackCommand(command string, duration time.Duration, success bool) {
	if !c.enabled {
		return
	}

	properties := map[string]interface{}{
		"command":  command,
		"duration": duration.Milliseconds(),
		"success":  success,
	}

	c.Track("command_executed", properties)

	// Add performance monitoring to Sentry
	if c.sentryEnabled {
		sentry.WithScope(func(scope *sentry.Scope) {
			scope.SetTag("command", command)
			scope.SetTag("success", fmt.Sprintf("%t", success))
			scope.SetExtra("duration_ms", duration.Milliseconds())

			// Create a transaction for performance monitoring
			ctx := sentry.SetHubOnContext(context.Background(), sentry.CurrentHub())
			transaction := sentry.StartTransaction(ctx, fmt.Sprintf("command.%s", command))
			transaction.SetTag("command", command)
			transaction.SetTag("success", fmt.Sprintf("%t", success))
			transaction.SetData("duration_ms", duration.Milliseconds())

			if !success {
				transaction.Status = sentry.SpanStatusInternalError
				sentry.CaptureMessage(fmt.Sprintf("Command failed: %s", command))
			} else {
				transaction.Status = sentry.SpanStatusOK
			}
			transaction.Finish()
		})
	}
}

func (c *Client) Close() {
	if c.segmentClient != nil {
		// Flush any pending events before closing
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Flushing Segment events before close\n")
		}
		c.segmentClient.Close()
	}
	if c.sentryEnabled {
		if os.Getenv("RUNE_DEBUG") == "true" {
			fmt.Printf("DEBUG: Flushing Sentry events before close\n")
		}
		sentry.Flush(5 * time.Second) // Increased timeout for better reliability
	}
}

// StartTransaction starts a Sentry transaction for performance monitoring
func (c *Client) StartTransaction(name, operation string) *sentry.Span {
	if !c.sentryEnabled {
		return nil
	}

	ctx := sentry.SetHubOnContext(context.Background(), sentry.CurrentHub())
	return sentry.StartTransaction(ctx, name)
}

// CaptureException captures an exception with additional context
func (c *Client) CaptureException(err error, tags map[string]string, extra map[string]interface{}) {
	if !c.sentryEnabled {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		for key, value := range tags {
			scope.SetTag(key, value)
		}
		for key, value := range extra {
			scope.SetExtra(key, value)
		}
		sentry.CaptureException(err)
	})
}

// CaptureMessage captures a message with additional context
func (c *Client) CaptureMessage(message string, level sentry.Level, tags map[string]string) {
	if !c.sentryEnabled {
		return
	}

	sentry.WithScope(func(scope *sentry.Scope) {
		scope.SetLevel(level)
		for key, value := range tags {
			scope.SetTag(key, value)
		}
		sentry.CaptureMessage(message)
	})
}

func getUserID() string {
	// Try to get from config first
	cfg, err := config.Load()
	if err == nil && cfg.UserID != "" {
		return cfg.UserID
	}

	// Generate a new anonymous ID
	userID := generateAnonymousID()

	// Try to save it to config
	if cfg != nil {
		cfg.UserID = userID
		_ = config.SaveConfig(cfg) // Ignore errors
	}

	return userID
}

func generateAnonymousID() string {
	// Simple anonymous ID generation
	hostname, _ := os.Hostname()
	return fmt.Sprintf("anon_%s_%d", hostname, time.Now().Unix())
}

func generateSessionID() string {
	return fmt.Sprintf("session_%d", time.Now().UnixNano())
}

// StartCommand starts tracking a command execution for release health
func (c *Client) StartCommand(command string) {
	if !c.sentryEnabled {
		return
	}

	// Add breadcrumb for command start
	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Message:   fmt.Sprintf("Command started: %s", command),
		Category:  "command",
		Level:     sentry.LevelInfo,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"command": command,
			"action":  "start",
		},
	})
}

// EndCommand ends tracking a command execution
func (c *Client) EndCommand(command string, success bool, duration time.Duration) {
	if !c.sentryEnabled {
		return
	}

	// Add breadcrumb for command end
	level := sentry.LevelInfo
	if !success {
		level = sentry.LevelError
	}

	sentry.AddBreadcrumb(&sentry.Breadcrumb{
		Message:   fmt.Sprintf("Command %s: %s", map[bool]string{true: "completed", false: "failed"}[success], command),
		Category:  "command",
		Level:     level,
		Timestamp: time.Now(),
		Data: map[string]interface{}{
			"command":     command,
			"action":      "end",
			"success":     success,
			"duration_ms": duration.Milliseconds(),
		},
	})
}

// Build-time variables set via ldflags
var (
	segmentWriteKey string
	sentryDSN       string
)

// Build-time version variable set via ldflags
var version string

func getVersion() string {
	if version != "" {
		return version
	}
	return "dev"
}

func getEnvironment() string {
	if env := os.Getenv("RUNE_ENV"); env != "" {
		return env
	}
	return "production"
}

func getOSVersion() string {
	switch runtime.GOOS {
	case "darwin":
		return getMacOSVersion()
	case "linux":
		return getLinuxVersion()
	case "windows":
		return getWindowsVersion()
	default:
		return "unknown"
	}
}

func getMacOSVersion() string {
	// Simple implementation - you might want to use a more robust method
	return "unknown"
}

func getLinuxVersion() string {
	// Simple implementation - you might want to use a more robust method
	return "unknown"
}

func getWindowsVersion() string {
	// Simple implementation - you might want to use a more robust method
	return "unknown"
}

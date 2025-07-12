package commands

import (
	"fmt"
	"net/http"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/spf13/cobra"
)

var debugCmd = &cobra.Command{
	Use:   "debug",
	Short: "Debug and diagnostic commands",
	Long:  "Debug and diagnostic commands for troubleshooting Rune issues",
}

var debugTelemetryCmd = &cobra.Command{
	Use:   "telemetry",
	Short: "Debug telemetry configuration and connectivity",
	Long:  "Show telemetry configuration, test connectivity, and verify event delivery",
	RunE:  runDebugTelemetry,
}

var debugKeysCmd = &cobra.Command{
	Use:   "keys",
	Short: "Show telemetry key configuration (masked for security)",
	Long:  "Display telemetry API keys and DSNs with masking for security",
	RunE:  runDebugKeys,
}

func init() {
	rootCmd.AddCommand(debugCmd)
	debugCmd.AddCommand(debugTelemetryCmd)
	debugCmd.AddCommand(debugKeysCmd)

	// Wrap debug commands with telemetry
	telemetry.WrapCommand(debugTelemetryCmd, runDebugTelemetry)
	telemetry.WrapCommand(debugKeysCmd, runDebugKeys)
}

func runDebugTelemetry(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Rune Telemetry Debug Report")
	fmt.Println("=" + strings.Repeat("=", 40))

	// System Information
	fmt.Printf("\n📊 System Information:\n")
	fmt.Printf("  OS: %s\n", runtime.GOOS)
	fmt.Printf("  Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("  Go Version: %s\n", runtime.Version())

	// Environment Variables
	fmt.Printf("\n🌍 Environment Variables:\n")
	fmt.Printf("  RUNE_TELEMETRY_DISABLED: %s\n", getEnvOrDefault("RUNE_TELEMETRY_DISABLED", "not set"))
	fmt.Printf("  RUNE_DEBUG: %s\n", getEnvOrDefault("RUNE_DEBUG", "not set"))
	fmt.Printf("  RUNE_ENV: %s\n", getEnvOrDefault("RUNE_ENV", "not set"))
	fmt.Printf("  RUNE_SEGMENT_WRITE_KEY: %s\n", maskKey(os.Getenv("RUNE_SEGMENT_WRITE_KEY")))
	fmt.Printf("  RUNE_SENTRY_DSN: %s\n", maskDSN(os.Getenv("RUNE_SENTRY_DSN")))

	// Configuration File
	fmt.Printf("\n📄 Configuration:\n")
	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("  Config Load Error: %v\n", err)
	} else {
		fmt.Printf("  Config File Found: ✅\n")
		fmt.Printf("  Telemetry Enabled: %t\n", cfg.Integrations.Telemetry.Enabled)
		fmt.Printf("  Segment Key (config): %s\n", maskKey(cfg.Integrations.Telemetry.SegmentWriteKey))
		fmt.Printf("  Sentry DSN (config): %s\n", maskDSN(cfg.Integrations.Telemetry.SentryDSN))
		fmt.Printf("  User ID: %s\n", cfg.UserID)
	}

	// Build-time Keys (check if embedded)
	fmt.Printf("\n🔧 Build-time Configuration:\n")
	fmt.Printf("  Build-time keys embedded: %s\n", checkBuildTimeKeys())

	// Network Connectivity Tests
	fmt.Printf("\n🌐 Network Connectivity:\n")
	testConnectivity("Segment API", "https://api.segment.io/v1/track")
	testConnectivity("Sentry API", "https://sentry.io/api/")

	// Test Event Sending
	fmt.Printf("\n📡 Test Event Sending:\n")
	testEventSending()

	return nil
}

func runDebugKeys(cmd *cobra.Command, args []string) error {
	fmt.Println("🔑 Rune API Keys Debug")
	fmt.Println("=" + strings.Repeat("=", 30))

	// Environment Variables
	segmentEnv := os.Getenv("RUNE_SEGMENT_WRITE_KEY")
	sentryEnv := os.Getenv("RUNE_SENTRY_DSN")

	fmt.Printf("\n🌍 Environment Variables:\n")
	fmt.Printf("  RUNE_SEGMENT_WRITE_KEY: %s\n", maskKey(segmentEnv))
	fmt.Printf("  RUNE_SENTRY_DSN: %s\n", maskDSN(sentryEnv))

	// Configuration File
	cfg, err := config.Load()
	if err == nil {
		fmt.Printf("\n📄 Configuration File:\n")
		fmt.Printf("  Segment Key: %s\n", maskKey(cfg.Integrations.Telemetry.SegmentWriteKey))
		fmt.Printf("  Sentry DSN: %s\n", maskDSN(cfg.Integrations.Telemetry.SentryDSN))
	}

	// Final Resolution
	finalSegment := segmentEnv
	if finalSegment == "" && cfg != nil {
		finalSegment = cfg.Integrations.Telemetry.SegmentWriteKey
	}

	finalSentry := sentryEnv
	if finalSentry == "" && cfg != nil {
		finalSentry = cfg.Integrations.Telemetry.SentryDSN
	}

	fmt.Printf("\n✅ Final Resolution:\n")
	fmt.Printf("  Active Segment Key: %s\n", maskKey(finalSegment))
	fmt.Printf("  Active Sentry DSN: %s\n", maskDSN(finalSentry))

	// Validation
	fmt.Printf("\n🔍 Validation:\n")
	if finalSegment == "" {
		fmt.Printf("  ❌ No Segment key configured\n")
	} else {
		fmt.Printf("  ✅ Segment key configured (%d chars)\n", len(finalSegment))
	}

	if finalSentry == "" {
		fmt.Printf("  ❌ No Sentry DSN configured\n")
	} else {
		fmt.Printf("  ✅ Sentry DSN configured (%d chars)\n", len(finalSentry))
	}

	return nil
}

func getEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func maskKey(key string) string {
	if key == "" {
		return "not set"
	}
	if len(key) <= 8 {
		return strings.Repeat("*", len(key))
	}
	return key[:4] + strings.Repeat("*", len(key)-8) + key[len(key)-4:]
}

func maskDSN(dsn string) string {
	if dsn == "" {
		return "not set"
	}
	// For Sentry DSN format: https://public_key@sentry.io/project_id
	if strings.Contains(dsn, "@") {
		parts := strings.Split(dsn, "@")
		if len(parts) >= 2 {
			return maskKey(parts[0]) + "@" + parts[1]
		}
	}
	return maskKey(dsn)
}

func checkBuildTimeKeys() string {
	// This is a simple check - in a real implementation, you'd check if build-time
	// variables were properly injected during the build process
	// For now, we'll indicate if the binary likely has embedded keys
	return "checking..." // This would need actual implementation
}

func testConnectivity(name, url string) {
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(url)
	if err != nil {
		fmt.Printf("  %s: ❌ Failed (%v)\n", name, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode < 400 {
		fmt.Printf("  %s: ✅ Connected (HTTP %d)\n", name, resp.StatusCode)
	} else {
		fmt.Printf("  %s: ⚠️  HTTP %d\n", name, resp.StatusCode)
	}
}

func testEventSending() {
	fmt.Printf("  Sending test event...\n")

	// Send a test event
	telemetry.Track("debug_test_event", map[string]interface{}{
		"test":      true,
		"timestamp": time.Now().Unix(),
		"source":    "debug_command",
	})

	fmt.Printf("  ✅ Test event sent (check your analytics dashboard)\n")
	fmt.Printf("  💡 Enable RUNE_DEBUG=true for detailed telemetry logs\n")
}

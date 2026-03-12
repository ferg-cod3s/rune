package otellogger

import (
	"os"
	"testing"
)

func TestOtelLoggerConfig_Defaults(t *testing.T) {
	config := &OtelLoggerConfig{}

	if config.ServiceName != "" {
		t.Errorf("expected empty ServiceName, got %q", config.ServiceName)
	}
	if config.ServiceVersion != "" {
		t.Errorf("expected empty ServiceVersion, got %q", config.ServiceVersion)
	}
	if config.DisableSentry {
		t.Error("expected DisableSentry to be false by default")
	}
	if config.DisableOTLP {
		t.Error("expected DisableOTLP to be false by default")
	}
}

func TestInitialize_DisabledViaTelemetryEnv(t *testing.T) {
	// Set environment to disable telemetry
	t.Setenv("RUNE_TELEMETRY_DISABLED", "true")

	logger, err := Initialize(&OtelLoggerConfig{
		ServiceName: "test-rune",
		DisableOTLP: true,
	})
	if err != nil {
		t.Fatalf("Initialize should not fail when telemetry is disabled: %v", err)
	}
	if logger == nil {
		t.Fatal("expected non-nil logger even when disabled")
	}
	if logger.enabled {
		t.Error("expected logger to be disabled")
	}

	// Close should not error
	if err := logger.Close(); err != nil {
		t.Errorf("Close should not error when disabled: %v", err)
	}
}

func TestInitialize_NilConfig(t *testing.T) {
	t.Setenv("RUNE_TELEMETRY_DISABLED", "true")

	logger, err := Initialize(nil)
	if err != nil {
		t.Fatalf("Initialize with nil config should not fail when disabled: %v", err)
	}
	if logger == nil {
		t.Fatal("expected non-nil logger")
	}
}

func TestOtelLogger_LogEvent_Disabled(t *testing.T) {
	logger := &OtelLogger{enabled: false}
	// Should not panic when calling methods on a disabled logger
	logger.LogEvent(0, "test event", "key", "value")
}

func TestOtelLogger_LogError_Disabled(t *testing.T) {
	logger := &OtelLogger{enabled: false}
	// Should not panic when calling LogError on a disabled logger
	logger.LogError(os.ErrNotExist, "test error", "key", "value")
}

func TestOtelLogger_GetLogger_NilLogger(t *testing.T) {
	logger := &OtelLogger{enabled: true, logger: nil}
	slogger := logger.GetLogger()
	if slogger == nil {
		t.Error("GetLogger should return non-nil slog.Logger even when internal logger is nil")
	}
}

func TestGetGlobalLogger_BeforeInit(t *testing.T) {
	// Reset global state
	oldGlobal := globalLogger
	globalLogger = nil
	defer func() { globalLogger = oldGlobal }()

	logger := GetGlobalLogger()
	if logger != nil {
		t.Error("expected nil global logger before initialization")
	}
}

func TestCloseGlobal_NilLogger(t *testing.T) {
	oldGlobal := globalLogger
	globalLogger = nil
	defer func() { globalLogger = oldGlobal }()

	if err := CloseGlobal(); err != nil {
		t.Errorf("CloseGlobal should not error when global logger is nil: %v", err)
	}
}

func TestGetVersion(t *testing.T) {
	// Without env var
	t.Setenv("RUNE_VERSION", "")
	v := getVersion()
	if v != "dev" {
		t.Errorf("expected 'dev', got %q", v)
	}

	// With env var
	t.Setenv("RUNE_VERSION", "1.0.0")
	v = getVersion()
	if v != "1.0.0" {
		t.Errorf("expected '1.0.0', got %q", v)
	}
}

func TestGetEnvironment(t *testing.T) {
	t.Setenv("RUNE_ENV", "")
	t.Setenv("ENVIRONMENT", "")
	env := getEnvironment()
	if env != "production" {
		t.Errorf("expected 'production', got %q", env)
	}

	t.Setenv("RUNE_ENV", "staging")
	env = getEnvironment()
	if env != "staging" {
		t.Errorf("expected 'staging', got %q", env)
	}
}

func TestGetOSVersion(t *testing.T) {
	v := getOSVersion()
	// Should return a non-empty string
	if v == "" {
		t.Error("expected non-empty OS version string")
	}
}

package commands

import (
	"testing"
	"time"
)

func TestPackageExists(t *testing.T) {
	// Basic test to ensure package compiles
	t.Log("Commands package compiles successfully")
}

func TestRootCommandExists(t *testing.T) {
	if rootCmd == nil {
		t.Fatal("rootCmd should not be nil")
	}
	if rootCmd.Use != "rune" {
		t.Errorf("expected rootCmd.Use 'rune', got %q", rootCmd.Use)
	}
}

func TestSubcommandsRegistered(t *testing.T) {
	expectedCommands := []string{
		"start", "stop", "pause", "resume", "status",
		"report", "init", "config", "ritual", "migrate",
		"update", "completion", "logs", "debug", "test",
	}

	registered := make(map[string]bool)
	for _, cmd := range rootCmd.Commands() {
		registered[cmd.Name()] = true
	}

	for _, expected := range expectedCommands {
		if !registered[expected] {
			t.Errorf("expected command %q to be registered", expected)
		}
	}
}

func TestMaskTelemetryKeyForLogging(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty key", "", "[not configured]"},
		{"short key", "abc", "[configured]"},
		{"normal key", "abcdef123456", "abcd****3456"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := maskTelemetryKeyForLogging(tt.input)
			if result != tt.expected {
				t.Errorf("maskTelemetryKeyForLogging(%q) = %q, expected %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestCalendarMonthsBetween(t *testing.T) {
	tests := []struct {
		name     string
		from     time.Time
		to       time.Time
		expected int
	}{
		{
			name:     "same day",
			from:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "exactly 1 month",
			from:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "almost 1 month",
			from:     time.Date(2025, 1, 15, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 2, 14, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "february to march",
			from:     time.Date(2025, 2, 4, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 3, 4, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "cross year",
			from:     time.Date(2024, 11, 15, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 2, 15, 0, 0, 0, 0, time.UTC),
			expected: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calendarMonthsBetween(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("calendarMonthsBetween(%v, %v) = %d, expected %d", tt.from, tt.to, result, tt.expected)
			}
		})
	}
}

func TestCalendarYearsBetween(t *testing.T) {
	tests := []struct {
		name     string
		from     time.Time
		to       time.Time
		expected int
	}{
		{
			name:     "same year",
			from:     time.Date(2025, 3, 4, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 6, 4, 0, 0, 0, 0, time.UTC),
			expected: 0,
		},
		{
			name:     "exactly 1 year",
			from:     time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2025, 3, 4, 0, 0, 0, 0, time.UTC),
			expected: 1,
		},
		{
			name:     "2 years across leap year",
			from:     time.Date(2024, 3, 4, 0, 0, 0, 0, time.UTC),
			to:       time.Date(2026, 3, 4, 0, 0, 0, 0, time.UTC),
			expected: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := calendarYearsBetween(tt.from, tt.to)
			if result != tt.expected {
				t.Errorf("calendarYearsBetween(%v, %v) = %d, expected %d", tt.from, tt.to, result, tt.expected)
			}
		})
	}
}

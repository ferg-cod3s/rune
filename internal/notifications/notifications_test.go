package notifications

import (
	"testing"
	"time"
)

func TestNotificationManager_Creation(t *testing.T) {
	nm := NewNotificationManager(true)
	if nm == nil {
		t.Fatal("NewNotificationManager returned nil")
	}
	if !nm.enabled {
		t.Error("Expected notifications to be enabled")
	}

	nmDisabled := NewNotificationManager(false)
	if nmDisabled.enabled {
		t.Error("Expected notifications to be disabled")
	}
}

func TestNotificationManager_DisabledSend(t *testing.T) {
	nm := NewNotificationManager(false)

	notification := Notification{
		Title:    "Test",
		Message:  "Test message",
		Type:     Custom,
		Priority: Normal,
		Sound:    false,
	}

	// Should not return an error when disabled
	err := nm.Send(notification)
	if err != nil {
		t.Errorf("Expected no error when notifications disabled, got: %v", err)
	}
}

func TestFormatDuration(t *testing.T) {
	tests := []struct {
		duration time.Duration
		expected string
	}{
		{30 * time.Second, "30 seconds"},
		{90 * time.Second, "1 minutes"},
		{5 * time.Minute, "5 minutes"},
		{65 * time.Minute, "1 hours 5 minutes"},
		{2 * time.Hour, "2 hours"},
		{150 * time.Minute, "2 hours 30 minutes"},
	}

	for _, test := range tests {
		result := formatDuration(test.duration)
		if result != test.expected {
			t.Errorf("formatDuration(%v) = %q, expected %q", test.duration, result, test.expected)
		}
	}
}

func TestIsSupported(t *testing.T) {
	// This will depend on the platform the test is running on
	supported := IsSupported()

	// We can't assert a specific value since it depends on the platform,
	// but we can ensure the function doesn't panic
	t.Logf("Notifications supported on this platform: %v", supported)
}

func TestNotificationTypes(t *testing.T) {
	nm := NewNotificationManager(true)

	// Test break reminder
	err := nm.SendBreakReminder(45 * time.Minute)
	if err != nil {
		t.Logf("Break reminder error (may be expected on CI): %v", err)
	}

	// Test end of day reminder
	err = nm.SendEndOfDayReminder(7*time.Hour+30*time.Minute, 8.0)
	if err != nil {
		t.Logf("End of day reminder error (may be expected on CI): %v", err)
	}

	// Test session complete
	err = nm.SendSessionComplete(2*time.Hour, "test-project")
	if err != nil {
		t.Logf("Session complete error (may be expected on CI): %v", err)
	}

	// Test idle detected
	err = nm.SendIdleDetected(10 * time.Minute)
	if err != nil {
		t.Logf("Idle detected error (may be expected on CI): %v", err)
	}
}

func TestGetSoundName(t *testing.T) {
	nm := NewNotificationManager(true)

	tests := []struct {
		notification Notification
		expected     string
	}{
		{
			Notification{Sound: false, Priority: Normal},
			"",
		},
		{
			Notification{Sound: true, Priority: Critical},
			"Basso",
		},
		{
			Notification{Sound: true, Priority: High},
			"Ping",
		},
		{
			Notification{Sound: true, Priority: Normal},
			"default",
		},
	}

	for _, test := range tests {
		result := nm.getSoundName(test.notification)
		if result != test.expected {
			t.Errorf("getSoundName(%+v) = %q, expected %q", test.notification, result, test.expected)
		}
	}
}

func TestGetUrgencyLevel(t *testing.T) {
	nm := NewNotificationManager(true)

	tests := []struct {
		priority Priority
		expected string
	}{
		{Critical, "critical"},
		{High, "normal"},
		{Normal, "normal"},
		{Low, "low"},
	}

	for _, test := range tests {
		result := nm.getUrgencyLevel(test.priority)
		if result != test.expected {
			t.Errorf("getUrgencyLevel(%v) = %q, expected %q", test.priority, result, test.expected)
		}
	}
}

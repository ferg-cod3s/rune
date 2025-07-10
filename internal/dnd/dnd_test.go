package dnd

import (
	"testing"
	"time"

	"github.com/ferg-cod3s/rune/internal/notifications"
)

func TestPackageExists(t *testing.T) {
	// Basic test to ensure package compiles
	// More comprehensive tests should be added as functionality is implemented
	t.Log("DND package test placeholder")
}

func TestNewDNDManager(t *testing.T) {
	nm := notifications.NewNotificationManager(true)
	dndManager := NewDNDManager(nm)

	if dndManager == nil {
		t.Fatal("NewDNDManager returned nil")
	}

	if dndManager.notificationManager != nm {
		t.Error("DND manager should store the notification manager reference")
	}
}

func TestDNDManagerWithoutNotifications(t *testing.T) {
	dndManager := NewDNDManager(nil)

	if dndManager == nil {
		t.Fatal("NewDNDManager returned nil")
	}

	// These should not panic even with nil notification manager
	err := dndManager.SendBreakNotification(30 * time.Minute)
	if err != nil {
		t.Errorf("SendBreakNotification should not error with nil notification manager: %v", err)
	}

	err = dndManager.SendEndOfDayNotification(8*time.Hour, 8.0)
	if err != nil {
		t.Errorf("SendEndOfDayNotification should not error with nil notification manager: %v", err)
	}

	err = dndManager.SendSessionCompleteNotification(2*time.Hour, "test")
	if err != nil {
		t.Errorf("SendSessionCompleteNotification should not error with nil notification manager: %v", err)
	}

	err = dndManager.SendIdleNotification(10 * time.Minute)
	if err != nil {
		t.Errorf("SendIdleNotification should not error with nil notification manager: %v", err)
	}
}

func TestDNDManagerWithNotifications(t *testing.T) {
	nm := notifications.NewNotificationManager(false) // Disabled to avoid actual notifications in tests
	dndManager := NewDNDManager(nm)

	// These should not error (notifications are disabled)
	err := dndManager.SendBreakNotification(30 * time.Minute)
	if err != nil {
		t.Errorf("SendBreakNotification failed: %v", err)
	}

	err = dndManager.SendEndOfDayNotification(8*time.Hour, 8.0)
	if err != nil {
		t.Errorf("SendEndOfDayNotification failed: %v", err)
	}

	err = dndManager.SendSessionCompleteNotification(2*time.Hour, "test")
	if err != nil {
		t.Errorf("SendSessionCompleteNotification failed: %v", err)
	}

	err = dndManager.SendIdleNotification(10 * time.Minute)
	if err != nil {
		t.Errorf("SendIdleNotification failed: %v", err)
	}
}

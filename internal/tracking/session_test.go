package tracking

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionState_String(t *testing.T) {
	tests := []struct {
		state    SessionState
		expected string
	}{
		{StateStopped, "Stopped"},
		{StateRunning, "Running"},
		{StatePaused, "Paused"},
		{SessionState(999), "Unknown"},
	}

	for _, tt := range tests {
		t.Run(tt.expected, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.state.String())
		})
	}
}

func TestTracker_StartStopSession(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Test starting a session
	session, err := tracker.Start("test-project")
	require.NoError(t, err)
	assert.Equal(t, "test-project", session.Project)
	assert.Equal(t, StateRunning, session.State)
	assert.False(t, session.StartTime.IsZero())
	assert.Nil(t, session.EndTime)

	// Test getting current session
	current, err := tracker.GetCurrentSession()
	require.NoError(t, err)
	require.NotNil(t, current)
	assert.Equal(t, session.ID, current.ID)
	assert.Equal(t, StateRunning, current.State)

	// Test starting another session should fail
	_, err = tracker.Start("another-project")
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session already active")

	// Test stopping the session
	stoppedSession, err := tracker.Stop()
	require.NoError(t, err)
	assert.Equal(t, session.ID, stoppedSession.ID)
	assert.Equal(t, StateStopped, stoppedSession.State)
	assert.NotNil(t, stoppedSession.EndTime)
	assert.True(t, stoppedSession.Duration > 0)

	// Test getting current session after stop
	current, err = tracker.GetCurrentSession()
	require.NoError(t, err)
	assert.Nil(t, current)
}

func TestTracker_PauseResumeSession(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Start a session
	session, err := tracker.Start("test-project")
	require.NoError(t, err)

	// Wait a bit to accumulate some time
	time.Sleep(10 * time.Millisecond)

	// Test pausing
	pausedSession, err := tracker.Pause()
	require.NoError(t, err)
	assert.Equal(t, session.ID, pausedSession.ID)
	assert.Equal(t, StatePaused, pausedSession.State)
	assert.NotNil(t, pausedSession.PausedAt)

	// Test pausing already paused session should fail
	_, err = tracker.Pause()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session is not running")

	// Test resuming
	resumedSession, err := tracker.Resume()
	require.NoError(t, err)
	assert.Equal(t, session.ID, resumedSession.ID)
	assert.Equal(t, StateRunning, resumedSession.State)
	assert.Nil(t, resumedSession.PausedAt)

	// Test resuming already running session should fail
	_, err = tracker.Resume()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "session is not paused")

	// Clean up
	_, err = tracker.Stop()
	require.NoError(t, err)
}

func TestTracker_GetSessionDuration(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Test duration when no session
	duration, err := tracker.GetSessionDuration()
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), duration)

	// Start a session
	_, err = tracker.Start("test-project")
	require.NoError(t, err)

	// Wait a bit
	time.Sleep(10 * time.Millisecond)

	// Test duration while running
	duration, err = tracker.GetSessionDuration()
	require.NoError(t, err)
	assert.True(t, duration > 0)

	// Pause and test duration
	_, err = tracker.Pause()
	require.NoError(t, err)

	pausedDuration, err := tracker.GetSessionDuration()
	require.NoError(t, err)
	assert.True(t, pausedDuration > 0)
	// Duration should not increase significantly while paused
	assert.True(t, pausedDuration <= duration+time.Millisecond)

	// Resume and stop
	_, err = tracker.Resume()
	require.NoError(t, err)

	_, err = tracker.Stop()
	require.NoError(t, err)

	// Test duration after stop should return 0 since there's no current session
	duration, err = tracker.GetSessionDuration()
	require.NoError(t, err)
	assert.Equal(t, time.Duration(0), duration)
}

func TestTracker_ErrorCases(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Test stopping when no session
	_, err := tracker.Stop()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active session to stop")

	// Test pausing when no session
	_, err = tracker.Pause()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no active session to pause")

	// Test resuming when no session
	_, err = tracker.Resume()
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no session to resume")
}

func TestTracker_GetDailyTotal(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Start and stop a session
	_, err := tracker.Start("test-project")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	_, err = tracker.Stop()
	require.NoError(t, err)

	// Get daily total
	dailyTotal, err := tracker.GetDailyTotal()
	require.NoError(t, err)
	assert.True(t, dailyTotal > 0)
}

func TestTracker_GetWeeklyTotal(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Start and stop a session
	_, err := tracker.Start("test-project")
	require.NoError(t, err)

	time.Sleep(10 * time.Millisecond)

	_, err = tracker.Stop()
	require.NoError(t, err)

	// Get weekly total
	weeklyTotal, err := tracker.GetWeeklyTotal()
	require.NoError(t, err)
	assert.True(t, weeklyTotal > 0)
}

func TestTracker_GetSessionHistory(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Create multiple sessions
	for i := 0; i < 3; i++ {
		_, err := tracker.Start(fmt.Sprintf("project-%d", i))
		require.NoError(t, err)

		time.Sleep(5 * time.Millisecond)

		_, err = tracker.Stop()
		require.NoError(t, err)
	}

	// Get session history
	sessions, err := tracker.GetSessionHistory(2)
	require.NoError(t, err)
	assert.Len(t, sessions, 2)

	// Should be ordered by most recent first
	assert.Equal(t, "project-2", sessions[0].Project)
	assert.Equal(t, "project-1", sessions[1].Project)
}

func TestTracker_GetProjectStats(t *testing.T) {
	tracker := setupTestTracker(t)
	defer tracker.Close()

	// Create sessions for different projects
	projects := []string{"project-a", "project-b", "project-a"}
	for _, project := range projects {
		_, err := tracker.Start(project)
		require.NoError(t, err)

		time.Sleep(5 * time.Millisecond)

		_, err = tracker.Stop()
		require.NoError(t, err)
	}

	// Get project stats
	stats, err := tracker.GetProjectStats()
	require.NoError(t, err)

	assert.Contains(t, stats, "project-a")
	assert.Contains(t, stats, "project-b")
	assert.True(t, stats["project-a"] > stats["project-b"]) // project-a has 2 sessions
}

// setupTestTracker creates a tracker with a temporary database for testing
func setupTestTracker(t *testing.T) *Tracker {
	// Create a temporary directory for the test database
	tempDir := t.TempDir()

	// Temporarily change HOME to use temp directory
	originalHome := os.Getenv("HOME")
	os.Setenv("HOME", tempDir)
	t.Cleanup(func() {
		os.Setenv("HOME", originalHome)
	})

	tracker, err := NewTracker()
	require.NoError(t, err)

	return tracker
}

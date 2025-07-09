package commands

import (
	"fmt"

	"github.com/johnferguson/rune/internal/telemetry"
	"github.com/johnferguson/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var resumeCmd = &cobra.Command{
	Use:   "resume",
	Short: "Resume the paused work timer",
	Long: `Resume a previously paused work session timer.

This command will:
- Resume the paused work timer
- Optionally re-enable focus mode
- Continue tracking the current session`,
	RunE: runResume,
}

func init() {
	rootCmd.AddCommand(resumeCmd)

	// Wrap command with telemetry
	telemetry.WrapCommand(resumeCmd, runResume)
}

func runResume(cmd *cobra.Command, args []string) error {
	fmt.Println("▶️ Resuming work timer...")

	// Initialize tracker
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to initialize tracker: %w", err)
	}
	defer tracker.Close()

	// Resume the session
	session, err := tracker.Resume()
	if err != nil {
		telemetry.TrackError(err, "resume", map[string]interface{}{
			"step": "tracker_resume",
		})
		return fmt.Errorf("failed to resume session: %w", err)
	}

	// Track successful resume
	telemetry.Track("session_resumed", map[string]interface{}{
		"project": session.Project,
	})

	fmt.Println("✓ Timer resumed")
	fmt.Println("🎯 Back to work!")

	return nil
}

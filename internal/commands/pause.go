package commands

import (
	"fmt"

	"github.com/ferg-cod3s/rune/internal/telemetry"
	"github.com/ferg-cod3s/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var pauseCmd = &cobra.Command{
	Use:   "pause",
	Short: "Pause the current work timer",
	Long: `Pause the current work session timer.

This command will:
- Pause the active work timer
- Optionally disable focus mode
- Save the current session state`,
	RunE: runPause,
}

func init() {
	rootCmd.AddCommand(pauseCmd)

	// Wrap command with telemetry
	telemetry.WrapCommand(pauseCmd, runPause)
}

func runPause(cmd *cobra.Command, args []string) error {
	fmt.Println("⏸ Pausing work timer...")

	// Initialize tracker
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to initialize tracker: %w", err)
	}
	defer tracker.Close()

	// Pause the session
	session, err := tracker.Pause()
	if err != nil {
		telemetry.TrackError(err, "pause", map[string]interface{}{
			"step": "tracker_pause",
		})
		return fmt.Errorf("failed to pause session: %w", err)
	}

	// Track successful pause
	telemetry.Track("session_paused", map[string]interface{}{
		"project": session.Project,
	})

	fmt.Println("✓ Timer paused")
	fmt.Println("💡 Use 'rune resume' to continue your session")

	return nil
}

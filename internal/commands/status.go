package commands

import (
	"fmt"

	"github.com/johnferguson/rune/internal/dnd"
	"github.com/johnferguson/rune/internal/telemetry"
	"github.com/johnferguson/rune/internal/tracking"
	"github.com/spf13/cobra"
)

var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show current session status",
	Long: `Display the current status of your work session.

This command shows:
- Current timer state (running, paused, stopped)
- Active project (if detected)
- Session duration
- Today's total work time
- Focus mode status`,
	RunE: runStatus,
}

func init() {
	rootCmd.AddCommand(statusCmd)

	// Wrap command with telemetry
	telemetry.WrapCommand(statusCmd, runStatus)
}

func runStatus(cmd *cobra.Command, args []string) error {
	fmt.Println("📊 Current Session Status")
	fmt.Println("========================")
	fmt.Println()

	// Initialize tracker
	tracker, err := tracking.NewTracker()
	if err != nil {
		return fmt.Errorf("failed to initialize tracker: %w", err)
	}
	defer tracker.Close()

	// Get current session
	session, err := tracker.GetCurrentSession()
	if err != nil {
		return fmt.Errorf("failed to get current session: %w", err)
	}

	if session == nil {
		fmt.Println("Timer:        Stopped")
		fmt.Println("Project:      Not detected")
		fmt.Println("Session:      0h 0m")
	} else {
		duration, err := tracker.GetSessionDuration()
		if err != nil {
			return fmt.Errorf("failed to get session duration: %w", err)
		}

		fmt.Printf("Timer:        %s\n", session.State)
		fmt.Printf("Project:      %s\n", session.Project)
		fmt.Printf("Session:      %s\n", formatDuration(duration))
	}

	// Get daily total
	dailyTotal, err := tracker.GetDailyTotal()
	if err != nil {
		return fmt.Errorf("failed to get daily total: %w", err)
	}
	fmt.Printf("Today Total:  %s\n", formatDuration(dailyTotal))

	// Get idle status
	isIdle, err := tracker.IsIdle()
	if err != nil {
		fmt.Println("Idle Status:  Unknown (detection failed)")
	} else {
		if isIdle {
			idleTime, err := tracker.GetIdleTime()
			if err == nil {
				fmt.Printf("Idle Status:  Idle for %s\n", formatDuration(idleTime))
			} else {
				fmt.Println("Idle Status:  Idle")
			}
		} else {
			fmt.Println("Idle Status:  Active")
		}
	}

	// Check DND status
	dndManager := dnd.NewDNDManager()
	dndEnabled, err := dndManager.IsEnabled()
	if err != nil {
		fmt.Println("Focus Mode:   Unknown (detection failed)")
	} else {
		if dndEnabled {
			fmt.Println("Focus Mode:   Enabled")
		} else {
			fmt.Println("Focus Mode:   Disabled")
		}
	}

	return nil
}

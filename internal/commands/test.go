package commands

import (
	"fmt"
	"time"

	"github.com/ferg-cod3s/rune/internal/dnd"
	"github.com/ferg-cod3s/rune/internal/notifications"
	"github.com/spf13/cobra"
)

// testCmd represents the test command
var testCmd = &cobra.Command{
	Use:   "test",
	Short: "Test various Rune functionality",
	Long: `Test various Rune functionality including notifications, DND, and integrations.

This command helps verify that your system is properly configured and that
Rune can interact with your operating system as expected.`,
}

// testNotificationsCmd tests the notification system
var testNotificationsCmd = &cobra.Command{
	Use:   "notifications",
	Short: "Test the notification system",
	Long: `Test the notification system to ensure OS-level notifications are working.

This will send various types of test notifications to verify that:
- Basic notifications work
- Break reminders work
- End-of-day reminders work
- Session completion notifications work
- Idle detection notifications work`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🧪 Testing notification system...")

		// Check if notifications are supported
		if !notifications.IsSupported() {
			return fmt.Errorf("notifications are not supported on this platform")
		}

		// Create notification manager
		nm := notifications.NewNotificationManager(true)

		// Create DND manager with notifications
		dndManager := dnd.NewDNDManager(nm)

		fmt.Println("📱 Sending test notification...")
		if err := dndManager.TestNotifications(); err != nil {
			fmt.Printf("❌ Test notification failed: %v\n", err)
		} else {
			fmt.Println("✅ Test notification sent successfully")
		}

		// Wait a moment between notifications
		time.Sleep(2 * time.Second)

		fmt.Println("🧘 Testing break reminder...")
		if err := dndManager.SendBreakNotification(45 * time.Minute); err != nil {
			fmt.Printf("❌ Break reminder failed: %v\n", err)
		} else {
			fmt.Println("✅ Break reminder sent successfully")
		}

		time.Sleep(2 * time.Second)

		fmt.Println("🌅 Testing end-of-day reminder...")
		if err := dndManager.SendEndOfDayNotification(7*time.Hour+30*time.Minute, 8.0); err != nil {
			fmt.Printf("❌ End-of-day reminder failed: %v\n", err)
		} else {
			fmt.Println("✅ End-of-day reminder sent successfully")
		}

		time.Sleep(2 * time.Second)

		fmt.Println("✅ Testing session complete notification...")
		if err := dndManager.SendSessionCompleteNotification(2*time.Hour, "test-project"); err != nil {
			fmt.Printf("❌ Session complete notification failed: %v\n", err)
		} else {
			fmt.Println("✅ Session complete notification sent successfully")
		}

		time.Sleep(2 * time.Second)

		fmt.Println("💤 Testing idle detection notification...")
		if err := dndManager.SendIdleNotification(10 * time.Minute); err != nil {
			fmt.Printf("❌ Idle detection notification failed: %v\n", err)
		} else {
			fmt.Println("✅ Idle detection notification sent successfully")
		}

		fmt.Println("\n🎉 Notification testing complete!")
		fmt.Println("If you saw notifications appear on your screen, the system is working correctly.")
		fmt.Println("If not, check your system's notification settings and permissions.")

		return nil
	},
}

// testDNDCmd tests the Do Not Disturb functionality
var testDNDCmd = &cobra.Command{
	Use:   "dnd",
	Short: "Test Do Not Disturb functionality",
	Long: `Test Do Not Disturb functionality to ensure Rune can control your system's
focus mode and notification settings.

This will:
- Check if DND is currently enabled
- Test enabling DND
- Test disabling DND
- Check for required shortcuts (macOS)`,
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Println("🔕 Testing Do Not Disturb functionality...")

		// Create DND manager
		nm := notifications.NewNotificationManager(true)
		dndManager := dnd.NewDNDManager(nm)

		// Check current status
		fmt.Println("📊 Checking current DND status...")
		enabled, err := dndManager.IsEnabled()
		if err != nil {
			fmt.Printf("❌ Failed to check DND status: %v\n", err)
		} else {
			fmt.Printf("ℹ️  DND is currently: %s\n", map[bool]string{true: "enabled", false: "disabled"}[enabled])
		}

		// Test enabling DND
		fmt.Println("🔕 Testing DND enable...")
		if err := dndManager.Enable(); err != nil {
			fmt.Printf("❌ Failed to enable DND: %v\n", err)
		} else {
			fmt.Println("✅ DND enabled successfully")
		}

		// Wait a moment
		time.Sleep(3 * time.Second)

		// Test disabling DND
		fmt.Println("🔔 Testing DND disable...")
		if err := dndManager.Disable(); err != nil {
			fmt.Printf("❌ Failed to disable DND: %v\n", err)
		} else {
			fmt.Println("✅ DND disabled successfully")
		}

		// Check shortcuts setup (macOS only)
		fmt.Println("🔧 Checking shortcuts setup...")
		shortcutsOK, err := dndManager.CheckShortcutsSetup()
		if err != nil {
			fmt.Printf("❌ Failed to check shortcuts: %v\n", err)
		} else if shortcutsOK {
			fmt.Println("✅ Required shortcuts are properly configured")
		} else {
			fmt.Println("⚠️  Some shortcuts may need to be configured manually")
		}

		fmt.Println("\n🎉 DND testing complete!")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(testCmd)
	testCmd.AddCommand(testNotificationsCmd)
	testCmd.AddCommand(testDNDCmd)
}

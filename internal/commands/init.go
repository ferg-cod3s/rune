package commands

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ferg-cod3s/rune/internal/config"
	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Initialize Rune configuration",
	Long: `Initialize your Rune configuration with guided setup.

This command will create a configuration file at ~/.rune/config.yaml
and walk you through setting up your first rituals and work preferences.`,
	RunE: runInit,
}

var (
	guided bool
)

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().BoolVar(&guided, "guided", false, "Use interactive guided setup")
}

func runInit(cmd *cobra.Command, args []string) error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	// Create .rune directory if it doesn't exist
	runeDir := filepath.Dir(configPath)
	if err := os.MkdirAll(runeDir, 0755); err != nil {
		return fmt.Errorf("failed to create .rune directory: %w", err)
	}

	// Check if config already exists
	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if exists {
		fmt.Printf("⚠ Configuration already exists at %s\n", configPath)
		fmt.Println("Use 'rune config edit' to modify your existing configuration.")
		return nil
	}

	if guided {
		return runGuidedInit(configPath)
	}

	return createDefaultConfig(configPath)
}

func runGuidedInit(configPath string) error {
	// Special ceremonial runic logo for initialization
	fmt.Println("|~\\  |\\      |   |\\    /|")
	fmt.Println("|  \\ | \\   \\ |   | \\  / |")
	fmt.Println("|  / |  \\   \\|   |  \\/  |")
	fmt.Println("|_/  |   |   |\\  |      |")
	fmt.Println("| \\  |   |   | \\ |      |")
	fmt.Println("|  \\ |   |   |   |      |")
	fmt.Println()
	fmt.Println("🔮 Ancient runes awaken... Welcome to Rune!")
	fmt.Println("Let's cast your daily rituals and bind your workflow.")
	fmt.Println()

	// Ask for telemetry opt-in
	telemetryEnabled := promptTelemetryOptIn()

	if err := createDefaultConfigWithTelemetry(configPath, telemetryEnabled); err != nil {
		return err
	}

	fmt.Printf("✓ Configuration created at %s\n", configPath)
	if telemetryEnabled {
		fmt.Println("✓ Telemetry enabled - helping improve Rune")
	} else {
		fmt.Println("✓ Telemetry disabled - fully private usage")
	}
	fmt.Println("✓ Ready to begin your ritual automation")
	fmt.Println()
	fmt.Println("Try 'rune start' to begin your workday!")

	return nil
}

func createDefaultConfig(configPath string) error {
	return createDefaultConfigWithTelemetry(configPath, false)
}

func createDefaultConfigWithTelemetry(configPath string, telemetryEnabled bool) error {
	telemetryConfig := "false"
	if telemetryEnabled {
		telemetryConfig = "true"
	}

	defaultConfig := fmt.Sprintf(`version: 1
settings:
  work_hours: 8.0
  break_interval: 50m
  idle_threshold: 10m

projects:
  - name: "default"
    detect: ["git:.*", "dir:~/"]

rituals:
  start:
    global:
      - name: "Welcome ritual"
        command: "echo 'Starting your workday...'"
  stop:
    global:
      - name: "Farewell ritual"
        command: "echo 'Ending your workday...'"

integrations:
  git:
    enabled: true
    auto_detect_project: true
  telemetry:
    enabled: %s
`, telemetryConfig)

	if err := os.WriteFile(configPath, []byte(defaultConfig), 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func promptTelemetryOptIn() bool {
	fmt.Println("📊 Help improve Rune")
	fmt.Println()
	fmt.Println("Rune can collect anonymous usage data and error reports to help")
	fmt.Println("improve the tool. This includes:")
	fmt.Println("  • Command usage patterns (no personal data)")
	fmt.Println("  • Error reports and crash logs")
	fmt.Println("  • Performance metrics")
	fmt.Println()
	fmt.Println("All data is anonymous and helps make Rune better for everyone.")
	fmt.Println("You can change this setting anytime with 'rune config edit'")
	fmt.Println()

	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enable telemetry? [y/N]: ")
		response, err := reader.ReadString('\n')
		if err != nil {
			return false
		}

		response = strings.TrimSpace(strings.ToLower(response))
		switch response {
		case "y", "yes":
			return true
		case "n", "no", "":
			return false
		default:
			fmt.Println("Please answer 'y' for yes or 'n' for no.")
		}
	}
}

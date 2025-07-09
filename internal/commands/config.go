package commands

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/johnferguson/rune/internal/config"
	"github.com/spf13/cobra"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Rune configuration",
	Long: `Manage your Rune configuration file.

This command provides subcommands to edit, validate, and manage
your Rune configuration.`,
}

var configEditCmd = &cobra.Command{
	Use:   "edit",
	Short: "Edit the configuration file",
	Long:  `Open the configuration file in your default editor.`,
	RunE:  runConfigEdit,
}

var configValidateCmd = &cobra.Command{
	Use:   "validate",
	Short: "Validate the configuration file",
	Long:  `Validate the syntax and content of your configuration file.`,
	RunE:  runConfigValidate,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show the current configuration",
	Long:  `Display the current configuration with resolved values.`,
	RunE:  runConfigShow,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configEditCmd)
	configCmd.AddCommand(configValidateCmd)
	configCmd.AddCommand(configShowCmd)
}

func runConfigEdit(cmd *cobra.Command, args []string) error {
	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("⚠ Configuration file does not exist.")
		fmt.Println("Run 'rune init' to create a new configuration.")
		return nil
	}

	editor := os.Getenv("EDITOR")
	if editor == "" {
		editor = "vi" // fallback to vi
	}

	cmd_exec := exec.Command(editor, configPath)
	cmd_exec.Stdin = os.Stdin
	cmd_exec.Stdout = os.Stdout
	cmd_exec.Stderr = os.Stderr

	return cmd_exec.Run()
}

func runConfigValidate(cmd *cobra.Command, args []string) error {
	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("⚠ Configuration file does not exist.")
		fmt.Println("Run 'rune init' to create a new configuration.")
		return nil
	}

	cfg, err := config.Load()
	if err != nil {
		fmt.Printf("❌ Configuration validation failed: %v\n", err)
		return nil // Don't return error to avoid double error message
	}

	fmt.Println("✅ Configuration is valid!")
	fmt.Printf("   Version: %d\n", cfg.Version)
	fmt.Printf("   Projects: %d\n", len(cfg.Projects))
	fmt.Printf("   Work hours: %.1f\n", cfg.Settings.WorkHours)

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	exists, err := config.Exists()
	if err != nil {
		return err
	}
	if !exists {
		fmt.Println("⚠ Configuration file does not exist.")
		fmt.Println("Run 'rune init' to create a new configuration.")
		return nil
	}

	configPath, err := config.GetConfigPath()
	if err != nil {
		return err
	}

	content, err := os.ReadFile(configPath)
	if err != nil {
		return fmt.Errorf("failed to read config file: %w", err)
	}

	fmt.Printf("Configuration file: %s\n", configPath)
	fmt.Println("=" + string(make([]byte, len(configPath)+20)))
	fmt.Print(string(content))

	return nil
}

package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
		errMsg  string
	}{
		{
			name: "valid config",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Projects: []Project{
					{Name: "test", Detect: []string{"git:test"}},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid version",
			config: Config{
				Version: 2,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "unsupported config version",
		},
		{
			name: "invalid work hours - zero",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "work_hours must be between 0 and 24",
		},
		{
			name: "invalid work hours - too high",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     25,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "work_hours must be between 0 and 24",
		},
		{
			name: "invalid break interval",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: -1 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "break_interval must be positive",
		},
		{
			name: "invalid idle threshold",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: -1 * time.Minute,
				},
			},
			wantErr: true,
			errMsg:  "idle_threshold must be positive",
		},
		{
			name: "project with empty name",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Projects: []Project{
					{Name: "", Detect: []string{"git:test"}},
				},
			},
			wantErr: true,
			errMsg:  "name cannot be empty",
		},
		{
			name: "project with empty detect patterns",
			config: Config{
				Version: 1,
				Settings: Settings{
					WorkHours:     8.0,
					BreakInterval: 50 * time.Minute,
					IdleThreshold: 10 * time.Minute,
				},
				Projects: []Project{
					{Name: "test", Detect: []string{}},
				},
			},
			wantErr: true,
			errMsg:  "detect patterns cannot be empty",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestGetConfigPath(t *testing.T) {
	path, err := GetConfigPath()
	require.NoError(t, err)

	home, err := os.UserHomeDir()
	require.NoError(t, err)

	expected := filepath.Join(home, ".rune", "config.yaml")
	assert.Equal(t, expected, path)
}

func TestExists(t *testing.T) {
	// Test when config doesn't exist
	// We'll use a temporary directory to avoid interfering with real config
	originalHome := os.Getenv("HOME")
	tempDir := t.TempDir()
	os.Setenv("HOME", tempDir)
	defer os.Setenv("HOME", originalHome)

	exists, err := Exists()
	require.NoError(t, err)
	assert.False(t, exists)

	// Create config file
	runeDir := filepath.Join(tempDir, ".rune")
	err = os.MkdirAll(runeDir, 0755)
	require.NoError(t, err)

	configPath := filepath.Join(runeDir, "config.yaml")
	err = os.WriteFile(configPath, []byte("version: 1"), 0644)
	require.NoError(t, err)

	// Test when config exists
	exists, err = Exists()
	require.NoError(t, err)
	assert.True(t, exists)
}

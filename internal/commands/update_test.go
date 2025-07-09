package commands

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsNewerVersion(t *testing.T) {
	tests := []struct {
		name     string
		latest   string
		current  string
		expected bool
	}{
		{
			name:     "newer version available",
			latest:   "v1.2.0",
			current:  "v1.1.0",
			expected: true,
		},
		{
			name:     "same version",
			latest:   "v1.1.0",
			current:  "v1.1.0",
			expected: false,
		},
		{
			name:     "same version without v prefix",
			latest:   "1.1.0",
			current:  "1.1.0",
			expected: false,
		},
		{
			name:     "mixed v prefix",
			latest:   "v1.2.0",
			current:  "1.2.0",
			expected: false,
		},
		{
			name:     "different versions",
			latest:   "v2.0.0",
			current:  "v1.9.9",
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := isNewerVersion(tt.latest, tt.current)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestFindAssetForPlatform(t *testing.T) {
	release := &GitHubRelease{
		TagName: "v1.0.0",
		Assets: []struct {
			Name               string `json:"name"`
			BrowserDownloadURL string `json:"browser_download_url"`
		}{
			{
				Name:               "rune_Darwin_x86_64.tar.gz",
				BrowserDownloadURL: "https://github.com/ferg-cod3s/rune/releases/download/v1.0.0/rune_Darwin_x86_64.tar.gz",
			},
			{
				Name:               "rune_Darwin_arm64.tar.gz",
				BrowserDownloadURL: "https://github.com/ferg-cod3s/rune/releases/download/v1.0.0/rune_Darwin_arm64.tar.gz",
			},
			{
				Name:               "rune_Linux_x86_64.tar.gz",
				BrowserDownloadURL: "https://github.com/ferg-cod3s/rune/releases/download/v1.0.0/rune_Linux_x86_64.tar.gz",
			},
			{
				Name:               "rune_Windows_x86_64.zip",
				BrowserDownloadURL: "https://github.com/ferg-cod3s/rune/releases/download/v1.0.0/rune_Windows_x86_64.zip",
			},
		},
	}

	tests := []struct {
		name        string
		expectedURL string
		shouldError bool
	}{
		{
			name:        "find asset for current platform",
			expectedURL: "", // Will be set based on runtime.GOOS/GOARCH
			shouldError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			url, err := findAssetForPlatform(release)

			if tt.shouldError {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, url)
				assert.Contains(t, url, "github.com")
			}
		})
	}
}

func TestUpdateCommand(t *testing.T) {
	// Test that the command is properly registered
	cmd := rootCmd
	updateCmd := cmd.Commands()

	var foundUpdate bool
	for _, subCmd := range updateCmd {
		if subCmd.Name() == "update" {
			foundUpdate = true
			break
		}
	}

	assert.True(t, foundUpdate, "update command should be registered")
}

func TestUpdateCommandFlags(t *testing.T) {
	// Test that the --check flag is properly configured
	flag := updateCmd.Flags().Lookup("check")
	assert.NotNil(t, flag, "check flag should exist")
	assert.Equal(t, "bool", flag.Value.Type(), "check flag should be boolean")
	assert.Equal(t, "false", flag.DefValue, "check flag should default to false")
}

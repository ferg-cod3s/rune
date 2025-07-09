package tracking

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProjectDetector_DetectProject(t *testing.T) {
	detector := NewProjectDetector()

	// Create a temporary directory for testing
	tempDir := t.TempDir()
	originalCwd, _ := os.Getwd()
	defer func() { _ = os.Chdir(originalCwd) }()

	t.Run("detects from package.json", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "node-project")
		require.NoError(t, os.MkdirAll(projectDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(projectDir, "package.json"), []byte(`{"name": "test-project"}`), 0644))

		require.NoError(t, os.Chdir(projectDir))
		project := detector.DetectProject()
		assert.Equal(t, "test-project", project)
	})

	t.Run("detects from go.mod", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "go-project")
		require.NoError(t, os.MkdirAll(projectDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(projectDir, "go.mod"), []byte("module test-project"), 0644))

		require.NoError(t, os.Chdir(projectDir))
		project := detector.DetectProject()
		assert.Equal(t, "test-project", project)
	})

	t.Run("detects from Cargo.toml", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "rust-project")
		require.NoError(t, os.MkdirAll(projectDir, 0755))
		require.NoError(t, os.WriteFile(filepath.Join(projectDir, "Cargo.toml"), []byte("[package]\nname = \"test-project\""), 0644))

		require.NoError(t, os.Chdir(projectDir))
		project := detector.DetectProject()
		assert.Equal(t, "test-project", project)
	})

	t.Run("detects from git repository", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "git-project")
		gitDir := filepath.Join(projectDir, ".git")
		require.NoError(t, os.MkdirAll(gitDir, 0755))

		require.NoError(t, os.Chdir(projectDir))
		project := detector.DetectProject()
		assert.Equal(t, "git-project", project)
	})

	t.Run("falls back to directory name", func(t *testing.T) {
		projectDir := filepath.Join(tempDir, "fallback-project")
		require.NoError(t, os.MkdirAll(projectDir, 0755))

		require.NoError(t, os.Chdir(projectDir))
		project := detector.DetectProject()
		assert.Equal(t, "fallback-project", project)
	})
}

func TestProjectDetector_SanitizeProjectName(t *testing.T) {
	detector := NewProjectDetector()

	tests := []struct {
		input    string
		expected string
	}{
		{"github.com/user/repo", "user-repo"},
		{"my-project.git", "my-project"},
		{"Project With Spaces", "project-with-spaces"},
		{"path/to/project", "path-to-project"},
		{"path\\to\\project", "path-to-project"},
		{"UPPERCASE", "uppercase"},
		{"", "default"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := detector.SanitizeProjectName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestProjectDetector_hasFile(t *testing.T) {
	detector := NewProjectDetector()
	tempDir := t.TempDir()

	// Create a test file
	testFile := filepath.Join(tempDir, "test.txt")
	require.NoError(t, os.WriteFile(testFile, []byte("test"), 0644))

	assert.True(t, detector.hasFile(tempDir, "test.txt"))
	assert.False(t, detector.hasFile(tempDir, "nonexistent.txt"))
}

func TestProjectDetector_isGitRepo(t *testing.T) {
	detector := NewProjectDetector()
	tempDir := t.TempDir()

	// Not a git repo initially
	assert.False(t, detector.isGitRepo(tempDir))

	// Create .git directory
	gitDir := filepath.Join(tempDir, ".git")
	require.NoError(t, os.MkdirAll(gitDir, 0755))

	// Now it should be detected as a git repo
	assert.True(t, detector.isGitRepo(tempDir))
}

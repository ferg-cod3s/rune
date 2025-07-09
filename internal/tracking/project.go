package tracking

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
)

// ProjectDetector handles automatic project detection
type ProjectDetector struct{}

// NewProjectDetector creates a new project detector
func NewProjectDetector() *ProjectDetector {
	return &ProjectDetector{}
}

// DetectProject attempts to detect the current project based on working directory
func (pd *ProjectDetector) DetectProject() string {
	cwd, err := os.Getwd()
	if err != nil {
		return "default"
	}

	// Check for common project indicators
	if pd.hasFile(cwd, "package.json") {
		return pd.getProjectNameFromPackageJSON(cwd)
	}

	if pd.hasFile(cwd, "go.mod") {
		return pd.getProjectNameFromGoMod(cwd)
	}

	if pd.hasFile(cwd, "Cargo.toml") {
		return pd.getProjectNameFromCargoToml(cwd)
	}

	if pd.hasFile(cwd, "pyproject.toml") || pd.hasFile(cwd, "setup.py") {
		return pd.getProjectNameFromPython(cwd)
	}

	// Check for git repository
	if pd.isGitRepo(cwd) {
		return pd.getProjectNameFromGit(cwd)
	}

	// Fall back to directory name
	return filepath.Base(cwd)
}

// hasFile checks if a file exists in the given directory
func (pd *ProjectDetector) hasFile(dir, filename string) bool {
	_, err := os.Stat(filepath.Join(dir, filename))
	return err == nil
}

// isGitRepo checks if the directory is a git repository
func (pd *ProjectDetector) isGitRepo(dir string) bool {
	return pd.hasFile(dir, ".git") || pd.findGitRoot(dir) != ""
}

// findGitRoot finds the git root directory
func (pd *ProjectDetector) findGitRoot(dir string) string {
	for {
		if pd.hasFile(dir, ".git") {
			return dir
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}
	return ""
}

// getProjectNameFromPackageJSON extracts project name from package.json
func (pd *ProjectDetector) getProjectNameFromPackageJSON(dir string) string {
	packagePath := filepath.Join(dir, "package.json")
	data, err := os.ReadFile(packagePath)
	if err != nil {
		return filepath.Base(dir)
	}

	var pkg struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return filepath.Base(dir)
	}

	if pkg.Name != "" {
		return pkg.Name
	}
	return filepath.Base(dir)
}

// getProjectNameFromGoMod extracts project name from go.mod
func (pd *ProjectDetector) getProjectNameFromGoMod(dir string) string {
	goModPath := filepath.Join(dir, "go.mod")
	file, err := os.Open(goModPath)
	if err != nil {
		return filepath.Base(dir)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	moduleRegex := regexp.MustCompile(`^module\s+(.+)$`)

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if matches := moduleRegex.FindStringSubmatch(line); len(matches) > 1 {
			moduleName := matches[1]
			// Extract just the project name from the module path
			parts := strings.Split(moduleName, "/")
			return parts[len(parts)-1]
		}
	}

	return filepath.Base(dir)
}

// getProjectNameFromCargoToml extracts project name from Cargo.toml
func (pd *ProjectDetector) getProjectNameFromCargoToml(dir string) string {
	cargoPath := filepath.Join(dir, "Cargo.toml")
	file, err := os.Open(cargoPath)
	if err != nil {
		return filepath.Base(dir)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	nameRegex := regexp.MustCompile(`name\s*=\s*"([^"]+)"`)
	inPackageSection := false

	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Check if we're entering the [package] section
		if line == "[package]" {
			inPackageSection = true
			continue
		}

		// Check if we're entering a different section
		if strings.HasPrefix(line, "[") && line != "[package]" {
			inPackageSection = false
			continue
		}

		// Look for name field in package section
		if inPackageSection {
			if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return filepath.Base(dir)
}

// getProjectNameFromPython extracts project name from Python project files
func (pd *ProjectDetector) getProjectNameFromPython(dir string) string {
	// Try pyproject.toml first
	pyprojectPath := filepath.Join(dir, "pyproject.toml")
	if file, err := os.Open(pyprojectPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		nameRegex := regexp.MustCompile(`^name\s*=\s*"([^"]+)"`)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
				return matches[1]
			}
		}
	}

	// Fall back to setup.py parsing (basic)
	setupPath := filepath.Join(dir, "setup.py")
	if file, err := os.Open(setupPath); err == nil {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		nameRegex := regexp.MustCompile(`name\s*=\s*['"']([^'"]+)['"']`)

		for scanner.Scan() {
			line := strings.TrimSpace(scanner.Text())
			if matches := nameRegex.FindStringSubmatch(line); len(matches) > 1 {
				return matches[1]
			}
		}
	}

	return filepath.Base(dir)
}

// getProjectNameFromGit extracts project name from git repository
func (pd *ProjectDetector) getProjectNameFromGit(dir string) string {
	gitRoot := pd.findGitRoot(dir)
	if gitRoot == "" {
		return filepath.Base(dir)
	}

	// Try to get remote origin URL
	cmd := exec.Command("git", "remote", "get-url", "origin")
	cmd.Dir = gitRoot
	output, err := cmd.Output()
	if err != nil {
		return filepath.Base(gitRoot)
	}

	remoteURL := strings.TrimSpace(string(output))
	if remoteURL != "" {
		// Extract project name from git URL
		// Handle both SSH and HTTPS URLs
		if strings.Contains(remoteURL, "github.com") || strings.Contains(remoteURL, "gitlab.com") || strings.Contains(remoteURL, "bitbucket.org") {
			// Remove .git suffix
			remoteURL = strings.TrimSuffix(remoteURL, ".git")
			// Extract the last part of the path
			parts := strings.Split(remoteURL, "/")
			if len(parts) > 0 {
				return parts[len(parts)-1]
			}
		}
	}

	return filepath.Base(gitRoot)
}

// SanitizeProjectName cleans up project names
func (pd *ProjectDetector) SanitizeProjectName(name string) string {
	// Remove common prefixes/suffixes
	name = strings.TrimPrefix(name, "github.com/")
	name = strings.TrimSuffix(name, ".git")

	// Replace invalid characters
	name = strings.ReplaceAll(name, "/", "-")
	name = strings.ReplaceAll(name, "\\", "-")
	name = strings.ReplaceAll(name, " ", "-")

	// Convert to lowercase
	name = strings.ToLower(name)

	if name == "" {
		return "default"
	}

	return name
}

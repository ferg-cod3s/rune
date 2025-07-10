package commands

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/spf13/cobra"
)

type GitHubRelease struct {
	TagName string `json:"tag_name"`
	Name    string `json:"name"`
	Assets  []struct {
		Name               string `json:"name"`
		BrowserDownloadURL string `json:"browser_download_url"`
	} `json:"assets"`
}

var updateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update rune to the latest version",
	Long: `Update rune to the latest version from GitHub releases.

This command will:
- Check for the latest version on GitHub
- Download and install the update if a newer version is available
- Preserve your current configuration

Examples:
  rune update              # Check and update to latest version
  rune update --check      # Only check for updates without installing`,
	RunE: runUpdate,
}

var (
	checkOnly bool
)

func init() {
	rootCmd.AddCommand(updateCmd)
	updateCmd.Flags().BoolVar(&checkOnly, "check", false, "Only check for updates without installing")
}

func runUpdate(cmd *cobra.Command, args []string) error {
	fmt.Println("🔍 Checking for updates...")

	// Get current version
	currentVersion := version
	if currentVersion == "dev" {
		return fmt.Errorf("cannot update development builds - please install from releases")
	}

	// Fetch latest release info
	latestRelease, err := getLatestRelease()
	if err != nil {
		return fmt.Errorf("failed to check for updates: %w", err)
	}

	// Compare versions
	if !isNewerVersion(latestRelease.TagName, currentVersion) {
		fmt.Printf("✅ You're already running the latest version (%s)\n", currentVersion)
		return nil
	}

	fmt.Printf("🆕 New version available: %s (current: %s)\n", latestRelease.TagName, currentVersion)

	if checkOnly {
		fmt.Println("💡 Run 'rune update' to install the latest version")
		return nil
	}

	// Find appropriate asset for current platform
	assetURL, err := findAssetForPlatform(latestRelease)
	if err != nil {
		return fmt.Errorf("failed to find download for your platform: %w", err)
	}

	fmt.Printf("⬇️  Downloading %s...\n", latestRelease.TagName)

	// Download and install
	if err := downloadAndInstall(assetURL, latestRelease.TagName); err != nil {
		return fmt.Errorf("failed to install update: %w", err)
	}

	fmt.Printf("✅ Successfully updated to %s!\n", latestRelease.TagName)
	fmt.Println("💡 Run 'rune --version' to verify the update")

	return nil
}

func getLatestRelease() (*GitHubRelease, error) {
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Get("https://api.github.com/repos/ferg-cod3s/rune/releases/latest")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("GitHub API returned status %d", resp.StatusCode)
	}

	var release GitHubRelease
	if err := json.NewDecoder(resp.Body).Decode(&release); err != nil {
		return nil, err
	}

	return &release, nil
}

func isNewerVersion(latest, current string) bool {
	// Remove 'v' prefix if present
	latest = strings.TrimPrefix(latest, "v")
	current = strings.TrimPrefix(current, "v")

	// Simple string comparison for now - in production you'd want proper semver comparison
	return latest != current
}

func findAssetForPlatform(release *GitHubRelease) (string, error) {
	osName := runtime.GOOS
	arch := runtime.GOARCH

	// Map Go arch names to release asset names
	archMap := map[string]string{
		"amd64": "x86_64",
		"arm64": "arm64",
	}

	if mappedArch, ok := archMap[arch]; ok {
		arch = mappedArch
	}

	// Map Go OS names to release asset names
	osMap := map[string]string{
		"darwin":  "Darwin",
		"linux":   "Linux",
		"windows": "Windows",
	}

	if mappedOS, ok := osMap[osName]; ok {
		osName = mappedOS
	}

	// Look for matching asset
	expectedPattern := fmt.Sprintf("rune_%s_%s", osName, arch)

	for _, asset := range release.Assets {
		if strings.Contains(asset.Name, expectedPattern) {
			return asset.BrowserDownloadURL, nil
		}
	}

	return "", fmt.Errorf("no release found for %s/%s", osName, arch)
}

func downloadAndInstall(url, version string) error {
	// Create temp directory
	tempDir, err := os.MkdirTemp("", "rune-update-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tempDir)

	// Download archive
	client := &http.Client{Timeout: 5 * time.Minute}
	resp, err := client.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("download failed with status %d", resp.StatusCode)
	}

	// Determine archive type and filename
	archivePath := filepath.Join(tempDir, "rune-archive")
	isZip := strings.HasSuffix(url, ".zip")

	if isZip {
		archivePath += ".zip"
	} else {
		archivePath += ".tar.gz"
	}

	// Save archive
	archiveFile, err := os.Create(archivePath)
	if err != nil {
		return err
	}

	_, err = io.Copy(archiveFile, resp.Body)
	archiveFile.Close()
	if err != nil {
		return err
	}

	// Extract archive
	extractDir := filepath.Join(tempDir, "extracted")
	if err := os.MkdirAll(extractDir, 0755); err != nil {
		return err
	}

	if isZip {
		err = extractZip(archivePath, extractDir)
	} else {
		err = extractTarGz(archivePath, extractDir)
	}
	if err != nil {
		return err
	}

	// Find the rune binary
	binaryName := "rune"
	if runtime.GOOS == "windows" {
		binaryName = "rune.exe"
	}

	var newBinaryPath string
	err = filepath.Walk(extractDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.Name() == binaryName {
			newBinaryPath = path
			return filepath.SkipDir
		}
		return nil
	})
	if err != nil {
		return err
	}

	if newBinaryPath == "" {
		return fmt.Errorf("could not find %s binary in archive", binaryName)
	}

	// Get current executable path
	currentExe, err := os.Executable()
	if err != nil {
		return err
	}

	// Make new binary executable
	if err := os.Chmod(newBinaryPath, 0755); err != nil {
		return err
	}

	// Replace current binary
	// On Windows, we might need to rename the old file first
	if runtime.GOOS == "windows" {
		backupPath := currentExe + ".old"
		if err := os.Rename(currentExe, backupPath); err != nil {
			return err
		}
		if err := copyFile(newBinaryPath, currentExe); err != nil {
			// Try to restore backup
			if restoreErr := os.Rename(backupPath, currentExe); restoreErr != nil {
				// Log the restore error but return the original error
				fmt.Fprintf(os.Stderr, "Warning: failed to restore backup: %v\n", restoreErr)
			}
			return err
		}
		if err := os.Remove(backupPath); err != nil {
			// Log warning but don't fail the update
			fmt.Fprintf(os.Stderr, "Warning: failed to remove backup file: %v\n", err)
		}
	} else {
		if err := copyFile(newBinaryPath, currentExe); err != nil {
			return err
		}
	}

	return nil
}

func extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		path := filepath.Join(dest, f.Name)

		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.FileInfo().Mode()); err != nil {
				return err
			}
			continue
		}

		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}

		rc, err := f.Open()
		if err != nil {
			return err
		}

		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func extractTarGz(src, dest string) error {
	file, err := os.Open(src)
	if err != nil {
		return err
	}
	defer file.Close()

	gzr, err := gzip.NewReader(file)
	if err != nil {
		return err
	}
	defer gzr.Close()

	tr := tar.NewReader(gzr)

	for {
		header, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		path := filepath.Join(dest, header.Name)

		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.MkdirAll(path, 0755); err != nil {
				return err
			}
		case tar.TypeReg:
			if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
				return err
			}

			outFile, err := os.Create(path)
			if err != nil {
				return err
			}

			if _, err := io.Copy(outFile, tr); err != nil {
				outFile.Close()
				return err
			}
			outFile.Close()

			if err := os.Chmod(path, os.FileMode(header.Mode)); err != nil {
				return err
			}
		}
	}

	return nil
}

func copyFile(src, dst string) error {
	sourceFile, err := os.Open(src)
	if err != nil {
		return err
	}
	defer sourceFile.Close()

	destFile, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer destFile.Close()

	if _, err := io.Copy(destFile, sourceFile); err != nil {
		return err
	}

	// Copy permissions
	sourceInfo, err := os.Stat(src)
	if err != nil {
		return err
	}

	return os.Chmod(dst, sourceInfo.Mode())
}

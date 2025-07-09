package tracking

import (
	"fmt"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"
)

// IdleDetector handles idle time detection across platforms
type IdleDetector struct {
	threshold time.Duration
}

// NewIdleDetector creates a new idle detector with the given threshold
func NewIdleDetector(threshold time.Duration) *IdleDetector {
	return &IdleDetector{
		threshold: threshold,
	}
}

// GetIdleTime returns the current system idle time
func (id *IdleDetector) GetIdleTime() (time.Duration, error) {
	switch runtime.GOOS {
	case "darwin":
		return id.getIdleTimeMacOS()
	case "linux":
		return id.getIdleTimeLinux()
	case "windows":
		return id.getIdleTimeWindows()
	default:
		return 0, fmt.Errorf("idle detection not supported on %s", runtime.GOOS)
	}
}

// IsIdle returns true if the system has been idle longer than the threshold
func (id *IdleDetector) IsIdle() (bool, error) {
	idleTime, err := id.GetIdleTime()
	if err != nil {
		return false, err
	}
	return idleTime >= id.threshold, nil
}

// getIdleTimeMacOS gets idle time on macOS using ioreg
func (id *IdleDetector) getIdleTimeMacOS() (time.Duration, error) {
	cmd := exec.Command("ioreg", "-c", "IOHIDSystem")
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to run ioreg: %w", err)
	}

	lines := strings.Split(string(output), "\n")
	for _, line := range lines {
		if strings.Contains(line, "HIDIdleTime") {
			// Extract the idle time value
			parts := strings.Split(line, "=")
			if len(parts) < 2 {
				continue
			}

			valueStr := strings.TrimSpace(parts[1])
			// Remove any trailing characters and parse as int64
			valueStr = strings.TrimSpace(strings.Split(valueStr, " ")[0])

			idleNanos, err := strconv.ParseInt(valueStr, 10, 64)
			if err != nil {
				continue
			}

			// Convert from nanoseconds to duration
			return time.Duration(idleNanos), nil
		}
	}

	return 0, fmt.Errorf("could not find HIDIdleTime in ioreg output")
}

// getIdleTimeLinux gets idle time on Linux using xprintidle or similar
func (id *IdleDetector) getIdleTimeLinux() (time.Duration, error) {
	// Try xprintidle first (most common)
	cmd := exec.Command("xprintidle")
	output, err := cmd.Output()
	if err == nil {
		idleMs, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
		if err != nil {
			return 0, fmt.Errorf("failed to parse xprintidle output: %w", err)
		}
		return time.Duration(idleMs) * time.Millisecond, nil
	}

	// Try xssstate as fallback
	cmd = exec.Command("xssstate", "-i")
	output, err = cmd.Output()
	if err == nil {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if strings.Contains(line, "idle:") {
				parts := strings.Fields(line)
				if len(parts) >= 2 {
					idleMs, err := strconv.ParseInt(parts[1], 10, 64)
					if err == nil {
						return time.Duration(idleMs) * time.Millisecond, nil
					}
				}
			}
		}
	}

	// Try parsing /proc/uptime and /proc/stat as last resort
	return id.getIdleTimeLinuxProc()
}

// getIdleTimeLinuxProc gets idle time from /proc filesystem (less accurate)
func (id *IdleDetector) getIdleTimeLinuxProc() (time.Duration, error) {
	// This is a simplified approach - in practice, calculating true idle time
	// from /proc/stat is complex and not very accurate for user idle detection
	return 0, fmt.Errorf("xprintidle not available and /proc method not implemented")
}

// getIdleTimeWindows gets idle time on Windows using GetLastInputInfo
func (id *IdleDetector) getIdleTimeWindows() (time.Duration, error) {
	// Use PowerShell to call GetLastInputInfo
	script := `
Add-Type @'
using System;
using System.Diagnostics;
using System.Runtime.InteropServices;

public struct LASTINPUTINFO {
    public uint cbSize;
    public uint dwTime;
}

public class Win32 {
    [DllImport("user32.dll")]
    public static extern bool GetLastInputInfo(ref LASTINPUTINFO plii);
    
    [DllImport("kernel32.dll")]
    public static extern uint GetTickCount();
}
'@

$lastInputInfo = New-Object LASTINPUTINFO
$lastInputInfo.cbSize = [System.Runtime.InteropServices.Marshal]::SizeOf($lastInputInfo)
[Win32]::GetLastInputInfo([ref]$lastInputInfo)
$idleTime = [Win32]::GetTickCount() - $lastInputInfo.dwTime
Write-Output $idleTime
`

	cmd := exec.Command("powershell", "-Command", script)
	output, err := cmd.Output()
	if err != nil {
		return 0, fmt.Errorf("failed to get idle time on Windows: %w", err)
	}

	idleMs, err := strconv.ParseInt(strings.TrimSpace(string(output)), 10, 64)
	if err != nil {
		return 0, fmt.Errorf("failed to parse Windows idle time: %w", err)
	}

	return time.Duration(idleMs) * time.Millisecond, nil
}

// StartIdleMonitoring starts monitoring for idle state changes
func (id *IdleDetector) StartIdleMonitoring(onIdleStart, onIdleEnd func()) chan struct{} {
	stop := make(chan struct{})

	go func() {
		ticker := time.NewTicker(30 * time.Second) // Check every 30 seconds
		defer ticker.Stop()

		wasIdle := false

		for {
			select {
			case <-stop:
				return
			case <-ticker.C:
				isIdle, err := id.IsIdle()
				if err != nil {
					// Log error but continue monitoring
					continue
				}

				if isIdle && !wasIdle {
					// Just became idle
					if onIdleStart != nil {
						onIdleStart()
					}
					wasIdle = true
				} else if !isIdle && wasIdle {
					// Just became active
					if onIdleEnd != nil {
						onIdleEnd()
					}
					wasIdle = false
				}
			}
		}
	}()

	return stop
}

package commands

import (
	"fmt"
	"time"
)

// formatDuration formats a duration as "Xh Ym"
func formatDuration(d time.Duration) string {
	hours := int(d.Hours())
	minutes := int(d.Minutes()) % 60
	return fmt.Sprintf("%dh %dm", hours, minutes)
}

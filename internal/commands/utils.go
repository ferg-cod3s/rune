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

// formatRelativeTime formats a time relative to now (e.g., "2h 30m ago", "just now", "yesterday")
func formatRelativeTime(t time.Time) string {
	now := time.Now()
	diff := now.Sub(t)

	// Handle future times (shouldn't happen in normal usage, but just in case)
	if diff < 0 {
		return "in the future"
	}

	// Less than 1 minute ago
	if diff < time.Minute {
		return "just now"
	}

	// Less than 1 hour ago - show minutes
	if diff < time.Hour {
		minutes := int(diff.Minutes())
		if minutes == 1 {
			return "1m ago"
		}
		return fmt.Sprintf("%dm ago", minutes)
	}

	// Less than 24 hours ago - show hours and minutes
	if diff < 24*time.Hour {
		hours := int(diff.Hours())
		minutes := int(diff.Minutes()) % 60
		if hours == 1 && minutes == 0 {
			return "1h ago"
		} else if minutes == 0 {
			return fmt.Sprintf("%dh ago", hours)
		} else if hours == 1 {
			return fmt.Sprintf("1h %dm ago", minutes)
		}
		return fmt.Sprintf("%dh %dm ago", hours, minutes)
	}

	// Check if it's yesterday
	yesterday := now.AddDate(0, 0, -1)
	if t.Year() == yesterday.Year() && t.Month() == yesterday.Month() && t.Day() == yesterday.Day() {
		return "yesterday"
	}

	// Less than 7 days ago - show days
	if diff < 7*24*time.Hour {
		days := int(diff.Hours() / 24)
		if days == 1 {
			return "1 day ago"
		}
		return fmt.Sprintf("%d days ago", days)
	}

	// Use calendar-based month calculation to handle variable month lengths
	months := calendarMonthsBetween(t, now)
	if months > 0 && months < 12 {
		if months == 1 {
			return "1 month ago"
		}
		return fmt.Sprintf("%d months ago", months)
	}

	// Between 7 days and 1 month - show weeks
	if months == 0 {
		weeks := int(diff.Hours() / (24 * 7))
		if weeks == 1 {
			return "1 week ago"
		}
		return fmt.Sprintf("%d weeks ago", weeks)
	}

	// More than a year ago - show years
	years := calendarYearsBetween(t, now)
	if years <= 1 {
		return "1 year ago"
	}
	return fmt.Sprintf("%d years ago", years)
}

// calendarMonthsBetween returns the number of whole calendar months between two times.
func calendarMonthsBetween(from, to time.Time) int {
	years := to.Year() - from.Year()
	months := int(to.Month()) - int(from.Month())
	total := years*12 + months
	// Adjust if the day hasn't been reached yet in the current month
	if to.Day() < from.Day() {
		total--
	}
	if total < 0 {
		return 0
	}
	return total
}

// calendarYearsBetween returns the number of whole calendar years between two times.
func calendarYearsBetween(from, to time.Time) int {
	years := to.Year() - from.Year()
	// Use month/day comparison to handle leap years correctly
	if to.Month() < from.Month() || (to.Month() == from.Month() && to.Day() < from.Day()) {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}

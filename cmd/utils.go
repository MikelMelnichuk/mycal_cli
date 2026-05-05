package cmd

import (
	"strings"
	"time"
)

const YYYYMMDD = "2006-01-02"
const HHMM = "15:04"
const AllDescription = "Show all events (including past ones)"
const JsonDescription = "Output as JSON"
const AfterDescription = "Filter events after given time (e.g., 12:00)"

// IsValidTime validates if the input matches HH:MM or H:MM format.
func isValidTime(s string) bool {
	_, err := time.Parse("15:04", strings.TrimSpace(s))
	return err == nil
}

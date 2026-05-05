package cmd

import (
	"strings"
	"time"
)

const YYYYMMDD = "2006-01-02"
const HHMM = "15:04"

// IsValidTime validates if the input matches HH:MM or H:MM format.
func isValidTime(s string) bool {
	_, err := time.Parse("15:04", strings.TrimSpace(s))
	return err == nil
}

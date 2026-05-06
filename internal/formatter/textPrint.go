package formatter

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/MikelMelnichuk/mycal/internal/models"
)

func PrettyPrintSingleDay(events []models.Event) {
	if len(events) == 0 {
		fmt.Println("No events for this day.")
		return
	}

	// Determine column headers and widths
	headers := []string{"Title", "Start", "End", "ID"}

	// Start with header lengths
	colWidths := []int{len(headers[0]), len(headers[1]), len(headers[2]), len(headers[3])}

	// Collect row data and adjust column widths
	rows := make([][]string, len(events))
	for i, e := range events {
		row := []string{e.Title, e.Start, e.End, e.ID}
		rows[i] = row
		for j, cell := range row {
			// For Title, limit width to 40 characters to avoid overly wide tables
			if j == 1 && utf8.RuneCountInString(cell) > 40 {
				cell = truncate(cell, 40)
				rows[i][j] = cell
			}
			w := utf8.RuneCountInString(cell)
			if w > colWidths[j] {
				colWidths[j] = w
			}
		}
	}

	// Build format string for each row
	format := "│ %-*s │ %-*s │ %-*s │ %-*s │\n"

	// Day header
	day := events[0].Day
	fmt.Printf("\n📅 Events for %s:\n\n", day)

	// Top border
	topBorder := "┌" + strings.Repeat("─", colWidths[0]+2) + "┬" +
		strings.Repeat("─", colWidths[1]+2) + "┬" +
		strings.Repeat("─", colWidths[2]+2) + "┬" +
		strings.Repeat("─", colWidths[3]+2) + "┐"
	fmt.Println(topBorder)

	// Header row
	fmt.Printf(format, colWidths[0], headers[0], colWidths[1], headers[1], colWidths[2], headers[2], colWidths[3], headers[3])

	// Separator
	sep := "├" + strings.Repeat("─", colWidths[0]+2) + "┼" +
		strings.Repeat("─", colWidths[1]+2) + "┼" +
		strings.Repeat("─", colWidths[2]+2) + "┼" +
		strings.Repeat("─", colWidths[3]+2) + "┤"
	fmt.Println(sep)

	// Data rows
	for _, row := range rows {
		fmt.Printf(format, colWidths[0], row[0], colWidths[1], row[1], colWidths[2], row[2], colWidths[3], row[3])
	}

	// Bottom border
	bottomBorder := "└" + strings.Repeat("─", colWidths[0]+2) + "┴" +
		strings.Repeat("─", colWidths[1]+2) + "┴" +
		strings.Repeat("─", colWidths[2]+2) + "┴" +
		strings.Repeat("─", colWidths[3]+2) + "┘"
	fmt.Println(bottomBorder)
}

// truncate shortens a string to maxLen characters and adds "..."
func truncate(s string, maxLen int) string {
	if utf8.RuneCountInString(s) <= maxLen {
		return s
	}
	// Ensure we don't cut in the middle of a UTF-8 sequence
	runes := []rune(s)
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

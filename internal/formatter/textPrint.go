package formatter

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"github.com/MikelMelnichuk/mycal/internal/models"
)

// PrettyPrintWeek prints events grouped by day for a week.
// Events can be for any set of days (max 7, but not required).
func PrettyPrintWeek(events []models.Event) {
	if len(events) == 0 {
		fmt.Println("No events this week.")
		return
	}

	start := 0
	for i, e := range events {
		if e.Day != events[start].Day {
			PrettyPrintSingleDay(events[start:i], false)
			start = i
		}
	}
	PrettyPrintSingleDay(events[start:], false)
}

func PrettyPrintSingleDay(events []models.Event, verbose bool) {
	if len(events) == 0 {
		fmt.Println("No events for this day.")
		return
	}

	// Determine column headers and widths
	headers := []string{"Title", "Start", "End"}
	colWidths := []int{len(headers[0]), len(headers[1]), len(headers[2])}

	if verbose {
		headers = append(headers, "ID")
		colWidths = append(colWidths, len("ID"))
	}

	numCols := len(headers)

	// Collect row data and adjust column widths
	rows := make([][]string, len(events))
	for i, e := range events {
		row := []string{e.Title, e.Start, e.End}
		if verbose {
			row = append(row, e.ID)
		}
		rows[i] = row
		for j, cell := range row {
			// Truncate the Title if needed
			if j == 0 && utf8.RuneCountInString(cell) > 40 {
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
	format := strings.Repeat("│ %-*s ", numCols) + "│\n"

	// Helper to print a single row
	printRow := func(cells []string) {
		args := make([]any, 0, numCols*2)
		for i := 0; i < numCols; i++ {
			args = append(args, colWidths[i], cells[i])
		}
		fmt.Printf(format, args...)
	}

	// Build top border, separator, bottom border dynamically
	buildBorder := func(left, mid, right string) string {
		var parts []string
		for i, w := range colWidths {
			if i == 0 {
				parts = append(parts, left+strings.Repeat("─", w+2))
			} else {
				parts = append(parts, strings.Repeat("─", w+2))
			}
		}
		return strings.Join(parts, mid) + right
	}

	topBorder := buildBorder("┌", "┬", "┐")
	sep := buildBorder("├", "┼", "┤")
	bottomBorder := buildBorder("└", "┴", "┘")

	// Day header
	day := events[0].Day
	fmt.Printf("\n📅 Events for %s:\n\n", day)

	// Print table
	fmt.Println(topBorder)
	printRow(headers)
	fmt.Println(sep)
	for _, row := range rows {
		printRow(row)
	}
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

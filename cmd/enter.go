// cmd/enter.go
package cmd

import (
	"fmt"
	"strings"

	"github.com/MikelMelnichuk/mycal/internal/models"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

var (
	enterSelectFlag int
	enterPrintFlag  bool
)

// enterCmd represents the interactive "enter" command
var enterCmd = &cobra.Command{
	Use:   "enter",
	Short: "Interactively explore events on a given day",
	Long: `Display events for a specific date and allow interactive selection to view full details.
Supports natural language dates (e.g., "next monday", "today") or ISO format (YYYY-MM-DD).`,
	Example: `  mycal enter today
  mycal enter today --select 1 --print
  mycal enter tomorrow
  mycal enter next friday
  mycal enter 2026-06-15`,
	RunE: runEnter,
}

func init() {
	enterCmd.Flags().IntVar(&enterSelectFlag, "select", 0, "Select an event by index (1-based) and print its details (requires --print)")
	enterCmd.Flags().BoolVar(&enterPrintFlag, "print", false, "Print the selected event details without entering the TUI")
	enterCmd.Flags().String("after", "", "Show events starting after given time (format HH:MM)")
	rootCmd.AddCommand(enterCmd)
}

func runEnter(cmd *cobra.Command, args []string) error {
	// Parse flags
	after, _ := cmd.Flags().GetString("after")
	all := true
	if after != "" {
		if !isValidTime(after) {
			return fmt.Errorf("invalid --after format %q, expected HH:MM", after)
		}
		all = false
	}

	if len(args) == 0 {
		return fmt.Errorf("no date description provided")
	}
	dateDesc := strings.Join(args, " ")

	// Fetch events for the given day
	events, err := APIClient.GetDayEvents(dateDesc, all, after)
	if err != nil {
		return err
	}
	if len(events) == 0 {
		fmt.Println("No events found for that day.")
		return nil
	}

	// Non-interactive mode: --print and --select
	if enterPrintFlag {
		if enterSelectFlag < 1 || enterSelectFlag > len(events) {
			return fmt.Errorf("--select %d is out of range (1-%d)", enterSelectFlag, len(events))
		}
		selected := events[enterSelectFlag-1]
		printEventDetails(selected)
		return nil
	}

	// Interactive TUI mode
	p := tea.NewProgram(newEnterModel(events))
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}
	return nil
}

// printEventDetails prints a single event in a human-friendly format.
// You can replace this with a more sophisticated formatter if needed.
func printEventDetails(e models.Event) {
	fmt.Printf("Title:       %s\n", e.Title)
	fmt.Printf("Start:       %s\n", e.Start)
	fmt.Printf("End:         %s\n", e.End)
	if e.Day != "" {
		fmt.Printf("Day: %s\n", e.Day)
	}
	// Add any other fields your Event type contains
}

// -------------------------------------------------------------------
// Bubbletea model for interactive event selection and detail view
// -------------------------------------------------------------------

type modelState int

const (
	listState modelState = iota
	detailState
)

type enterModel struct {
	events      []models.Event
	state       modelState
	selectedIdx int // index of the selected event (0-based)
	width       int
	height      int
	listStyle   lipgloss.Style
	detailStyle lipgloss.Style
	helpStyle   lipgloss.Style
}

func newEnterModel(events []models.Event) enterModel {
	return enterModel{
		events:      events,
		state:       listState,
		selectedIdx: 0,
		listStyle:   lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()),
		detailStyle: lipgloss.NewStyle().Padding(1).Border(lipgloss.RoundedBorder()),
		helpStyle:   lipgloss.NewStyle().Foreground(lipgloss.Color("240")).Padding(0, 1),
	}
}

func (m enterModel) Init() tea.Cmd {
	return nil
}

func (m enterModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch m.state {
		case listState:
			switch msg.String() {
			case "q", "ctrl+c":
				return m, tea.Quit
			case "up", "k":
				if m.selectedIdx > 0 {
					m.selectedIdx--
				}
			case "down", "j":
				if m.selectedIdx < len(m.events)-1 {
					m.selectedIdx++
				}
			case "enter":
				m.state = detailState
			}
		case detailState:
			switch msg.String() {
			case "b", "left", "backspace":
				m.state = listState
			case "q", "esc", "ctrl+c":
				return m, tea.Quit
			}
		}
	}
	return m, nil
}

func (m enterModel) View() string {
	if m.state == listState {
		return m.renderList()
	}
	return m.renderDetail()
}

// Improved renderList for enterModel
func (m enterModel) renderList() string {
	const (
		idxWidth  = 4
		minTitleW = 20
		minTimeW  = 12
		minLocW   = 10
	)

	// Determine available width
	availableWidth := m.width - 10 // inner padding/border margins
	if availableWidth < 40 {
		availableWidth = 40 // fallback for tiny terminals
	}

	// Column widths: index fixed, time fixed, location flexible, title takes rest
	timeWidth := 25 // combined "HH:MM - HH:MM"
	dayWidth := 15
	titleWidth := availableWidth - idxWidth - timeWidth - dayWidth
	if titleWidth < minTitleW {
		// Not enough space: shrink location or time
		dayWidth = minLocW
		timeWidth = minTimeW
		titleWidth = availableWidth - idxWidth - timeWidth - dayWidth
		if titleWidth < minTitleW {
			titleWidth = minTitleW
		}
	}

	// Header style
	headerStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Padding(0, 1)

	// Row styles
	evenRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("252"))
	oddRowStyle := lipgloss.NewStyle().Foreground(lipgloss.Color("255"))
	selectedStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("205")).
		Bold(true).
		Background(lipgloss.Color("235"))

	// Separator line
	separator := strings.Repeat("─", availableWidth+4) // +4 for borders

	// Build header
	header := fmt.Sprintf("%-*s %-*s %-*s",
		idxWidth, "#",
		titleWidth, "TITLE",
		timeWidth, "TIME")
	header = headerStyle.Render(header)

	var rows []string
	for i, ev := range m.events {
		idx := fmt.Sprintf("%d.", i+1)
		title := truncate(ev.Title, titleWidth)
		timeStr := fmt.Sprintf("%s – %s", ev.Start, ev.End)

		line := fmt.Sprintf("%-*s %-*s %-*s",
			idxWidth, idx,
			titleWidth, title,
			timeWidth, timeStr)

		// Apply row style
		var styledLine string
		if i == m.selectedIdx {
			styledLine = selectedStyle.Render(line)
		} else {
			if i%2 == 0 {
				styledLine = evenRowStyle.Render(line)
			} else {
				styledLine = oddRowStyle.Render(line)
			}
		}
		rows = append(rows, styledLine)
	}

	// Combine everything into a bordered box
	content := lipgloss.JoinVertical(lipgloss.Top,
		header,
		separator,
		lipgloss.JoinVertical(lipgloss.Top, rows...),
	)

	listBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(m.width - 2).
		Render(content)

	help := m.helpStyle.Render("↑/↓ or j/k • enter to select • q to quit")
	return lipgloss.JoinVertical(lipgloss.Top, listBox, help)
}

// Helper to truncate strings with ellipsis
func truncate(s string, maxLen int) string {
	if maxLen < 3 {
		maxLen = 3
	}
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

// renderDetail shows event details in a structured, pretty box.
func (m enterModel) renderDetail() string {
	ev := m.events[m.selectedIdx]

	// Determine usable width (leave margins for borders/padding)
	innerWidth := m.width - 6 // border + inner padding
	if innerWidth < 40 {
		innerWidth = 40
	}

	// Label style (bold cyan)
	labelStyle := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color("39")).
		Width(14) // fixed label width for alignment

	// Value style (normal white, with wrapping)
	valueStyle := lipgloss.NewStyle().
		Foreground(lipgloss.Color("255")).
		Width(innerWidth - 16). // leave space for label
		Align(lipgloss.Left)

	// Helper to render a labeled field with automatic wrapping
	renderField := func(label, value string) string {
		if value == "" {
			return ""
		}
		wrapped := valueStyle.Render(lipgloss.NewStyle().Width(innerWidth - 16).Render(value))
		return labelStyle.Render(label+":") + " " + wrapped + "\n"
	}

	// Build content
	var b strings.Builder
	b.WriteString(renderField("Day", ev.Day))
	b.WriteString(renderField("Title", ev.Title))
	b.WriteString(renderField("Start", ev.Start))
	b.WriteString(renderField("End", ev.End))

	detailContent := b.String()
	if detailContent == "" {
		detailContent = "No details available for this event."
	}

	// Box styling – same as list view
	detailBox := lipgloss.NewStyle().
		Border(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color("240")).
		Padding(1, 2).
		Width(m.width - 2).
		Render(detailContent)

	// Help text
	help := m.helpStyle.Render("←/b back • q/esc quit")
	return lipgloss.JoinVertical(lipgloss.Top, detailBox, help)
}

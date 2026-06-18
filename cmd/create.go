// cmd/create.go
package cmd

import (
	"fmt"
	"time"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/spf13/cobra"
)

// createCmd represents the interactive "create" command
var createCmd = &cobra.Command{
	Use:   "create",
	Short: "Interactively create a new event",
	Long: `Launch a TUI to enter event details (title, start/end times, description, location)
and send a POST request to the backend to create the event.`,
	Example: `  mycal create`,
	RunE:    runCreate,
}

func init() {
	rootCmd.AddCommand(createCmd)
}

func runCreate(cmd *cobra.Command, args []string) error {
	p := tea.NewProgram(newCreateModel())
	if _, err := p.Run(); err != nil {
		return fmt.Errorf("failed to start TUI: %w", err)
	}
	return nil
}

// -------------------------------------------------------------------
// Bubbletea model for the interactive create form
// -------------------------------------------------------------------

const (
	fieldTitle = iota
	fieldStart
	fieldEnd
	fieldDescription
	fieldLocation
	numFields
)

// createModel holds the TUI state
type createModel struct {
	inputs     []textinput.Model
	focusIndex int
	width      int
	height     int
	err        error
	submitted  bool // true after successful creation
	loading    bool // true while waiting for server response
	formStyle  lipgloss.Style
	helpStyle  lipgloss.Style
	labelStyle lipgloss.Style
	errorStyle lipgloss.Style
	successMsg string
}

func newCreateModel() createModel {
	// Initialize text inputs with placeholders and validation
	inputs := make([]textinput.Model, numFields)

	inputs[fieldTitle] = textinput.New()
	inputs[fieldTitle].Placeholder = "Event title"
	inputs[fieldTitle].Focus()
	inputs[fieldTitle].CharLimit = 100

	default_start := time.Now().Round(10 * time.Minute)
	inputs[fieldStart] = textinput.New()
	inputs[fieldStart].SetValue(default_start.Format("2006-01-02 15:04"))
	inputs[fieldStart].CharLimit = 30

	inputs[fieldEnd] = textinput.New()
	default_end := default_start.Add(time.Hour)
	inputs[fieldEnd].SetValue(default_end.Format("2006-01-02 15:04"))
	inputs[fieldEnd].CharLimit = 30

	inputs[fieldLocation] = textinput.New()
	inputs[fieldLocation].Placeholder = "Location (optional)"
	inputs[fieldLocation].CharLimit = 100

	inputs[fieldDescription] = textinput.New()
	inputs[fieldDescription].Placeholder = "Event description (optional)"
	inputs[fieldDescription].CharLimit = 500

	m := createModel{
		inputs:     inputs,
		focusIndex: 0,
		formStyle: lipgloss.NewStyle().
			Border(lipgloss.RoundedBorder()).
			BorderForeground(lipgloss.Color("240")).
			Padding(1, 2),
		helpStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("240")).
			Padding(0, 1),
		labelStyle: lipgloss.NewStyle().
			Bold(true).
			Foreground(lipgloss.Color("39")).
			Width(14),
		errorStyle: lipgloss.NewStyle().
			Foreground(lipgloss.Color("196")),
	}

	return m
}

func (m createModel) Init() tea.Cmd {
	return textinput.Blink
}

// Custom messages for async submission
type errMsg struct{ err error }
type successMsg struct{}

func (m createModel) View() string {
	if m.submitted {
		content := lipgloss.JoinVertical(
			lipgloss.Center,
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("46")).
				Bold(true).
				Render(m.successMsg),
			lipgloss.NewStyle().
				Foreground(lipgloss.Color("240")).
				Render("Press any key to quit"),
		)
		return m.formStyle.Width(m.width - 2).Render(content)
	}
	// Build the form
	var rows []string
	labels := []string{"Title", "Start", "End", "Location", "Description"}

	for i := 0; i < numFields; i++ {
		label := m.labelStyle.Render(labels[i] + ":")
		input := m.inputs[i].View()
		// If the input is focused, highlight it
		if m.focusIndex == i {
			input = lipgloss.NewStyle().Background(lipgloss.Color("235")).Render(input)
		}
		rows = append(rows, lipgloss.JoinHorizontal(lipgloss.Top, label, input))
	}

	formContent := lipgloss.JoinVertical(lipgloss.Top, rows...)

	// Show loading indicator or error
	var status string
	if m.loading {
		status = lipgloss.NewStyle().Foreground(lipgloss.Color("214")).Render("Sending request...")
	} else if m.err != nil {
		status = m.errorStyle.Render("Error: " + m.err.Error())
	}

	// Help text
	help := m.helpStyle.Render("↑/↓ or tab to navigate • enter to submit (on last field) • ctrl+c to quit")

	// Assemble
	content := lipgloss.JoinVertical(lipgloss.Top,
		formContent,
		status,
		help,
	)

	return m.formStyle.Width(m.width - 2).Render(content)
}

func (m createModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if m.submitted || m.loading {
		// In these states, we only handle quit keys or custom messages
		switch msg := msg.(type) {
		case tea.KeyMsg:
			if msg.String() == "ctrl+c" || (m.submitted && msg.String() == "q") {
				return m, tea.Quit
			}
			if m.submitted {
				// Any key quits
				return m, tea.Quit
			}
		case errMsg:
			m.loading = false
			m.err = msg.err
			return m, nil
		case successMsg:
			m.loading = false
			m.submitted = true
			m.successMsg = "✓ Event created successfully! Press any key to quit."
			return m, tea.ClearScreen
		}
		return m, nil
	}

	// Normal editing state
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width = msg.Width
		m.height = msg.Height
		return m, nil

	case tea.KeyMsg:
		switch msg.String() {
		case "ctrl+c":
			return m, tea.Quit

		case "tab", "shift+tab", "up", "down":
			s := msg.String()
			if s == "up" || s == "shift+tab" {
				m.focusIndex--
				if m.focusIndex < 0 {
					m.focusIndex = numFields - 1
				}
			} else {
				m.focusIndex++
				if m.focusIndex >= numFields {
					m.focusIndex = 0
				}
			}
			for i := 0; i < numFields; i++ {
				if i == m.focusIndex {
					m.inputs[i].Focus()
				} else {
					m.inputs[i].Blur()
				}
			}
			return m, nil

		case "enter":
			if m.focusIndex == numFields-1 {
				// Submit command
				return m, m.submit()
			} else {
				m.focusIndex++
				if m.focusIndex >= numFields {
					m.focusIndex = 0
				}
				for i := 0; i < numFields; i++ {
					if i == m.focusIndex {
						m.inputs[i].Focus()
					} else {
						m.inputs[i].Blur()
					}
				}
				return m, nil
			}
		}
	}

	// Update the focused input
	var cmd tea.Cmd
	m.inputs[m.focusIndex], cmd = m.inputs[m.focusIndex].Update(msg)
	return m, cmd
}

// Note: The submit method returns tea.Cmd. It should trigger an async operation.
// We'll define it as a method that returns a tea.Cmd.
func (m *createModel) submit() tea.Cmd {
	// Validation and create the command
	title := m.inputs[fieldTitle].Value()
	start := m.inputs[fieldStart].Value()
	end := m.inputs[fieldEnd].Value()

	if title == "" || start == "" || end == "" {
		m.err = fmt.Errorf("title, start, and end are required")
		return nil
	}

	// Parse start and end times using flexible format
	startTime, err := parseFlexibleTime(m.inputs[fieldStart].Value())
	if err != nil {
		m.err = fmt.Errorf("start time: %w", err)
		return nil
	}
	endTime, err := parseFlexibleTime(m.inputs[fieldEnd].Value())
	if err != nil {
		m.err = fmt.Errorf("end time: %w", err)
		return nil
	}
	// Convert to RFC3339 for the API (UTC, with 'Z')
	startRFC := startTime.Format(time.RFC3339)
	endRFC := endTime.Format(time.RFC3339)

	m.loading = true
	m.err = nil

	// Return a command that executes the request and sends a message
	return func() tea.Msg {
		success, err := APIClient.CreateEvent(title,
			startRFC,
			endRFC,
			m.inputs[fieldDescription].Value(),
			m.inputs[fieldLocation].Value())
		if err != nil {
			return errMsg{err: err}
		}
		if !success {
			return errMsg{err: fmt.Errorf("server returned is_successful: false")}
		}
		return successMsg{}
	}
}

func parseFlexibleTime(s string) (time.Time, error) {
	// Layouts that include a timezone offset (e.g., Z, +02:00, -05:00)
	layoutsWithTZ := []string{
		time.RFC3339,               // "2006-01-02T15:04:05Z07:00"
		"2006-01-02T15:04:05.000Z", // with millis and Z
		"2006-01-02T15:04:05Z",     // with Z (no offset)
	}
	// Layouts that do NOT include a timezone – these are assumed to be in local time
	layoutsLocal := []string{
		"2006-01-02T15:04:05",     // ISO with T but no offset
		"2006-01-02T15:04",        // no seconds
		"2006-01-02 15:04:05.000", // space with millis
		"2006-01-02 15:04:05",     // space with seconds
		"2006-01-02 15:04",        // space, no seconds
		"2006-01-02",              // date only → midnight local
	}

	// Try layouts that contain a timezone first (use time.Parse to keep offset)
	for _, layout := range layoutsWithTZ {
		t, err := time.Parse(layout, s)
		if err == nil {
			return t, nil // location is whatever the input specified
		}
	}
	// Try local layouts with time.ParseInLocation (assigns local timezone)
	for _, layout := range layoutsLocal {
		t, err := time.ParseInLocation(layout, s, time.Local)
		if err == nil {
			return t, nil
		}
	}
	return time.Time{}, fmt.Errorf("unable to parse %q as a date/time", s)
}

## Recommended Language: **Go**

**Go is the best choice** for your `mycal` CLI tool. Here's why, plus complete implementation strategy:

### Why Go? 🎯

| Factor | Go ✅ | Python | Node.js | Rust |
|--------|-------|--------|---------|------|
| **Single binary** | Native, no deps | Needs Python | Needs Node | Native, complex |
| **Performance** | Fast HTTP calls | Good enough | Good enough | Overkill |
| **CLI ecosystem** | `cobra`, `urfave/cli` | `click`, `typer` | `commander` | `clap` |
| **Interactive TUI** | `bubbletea`, `charmbracelet` | `textual`, `rich` | `ink`, `blessed` | `ratatui` |
| **Cross-platform** | Excellent | Good | Good | Good |
| **HTTP client** | Built-in | `requests` | Built-in | `reqwest` |
| **Learning curve** | Low for CLI | Lowest | Medium | High |

**Go wins** because it produces a single binary you can drop anywhere—no Python/Node install needed.

## 🏗️ Implementation Strategy

```
mycal/
├── cmd/
│   ├── root.go         # mycal
│   ├── today.go        # mycal today
│   ├── tomorrow.go     # mycal tomorrow  
│   ├── date.go         # mycal date
│   ├── week.go         # mycal week
│   ├── nextweek.go     # mycal nextweek
│   ├── enter.go        # mycal enter
│   └── test.go         # mycal test
├── api/                # HTTP client + models
│   ├── client.go
│   └── models.go
├── date/               # Date parsing
│   └── parser.go
├── ui/                 # Output formatting
│   ├── list.go
│   └── interactive.go
└── main.go
```

## Core Dependencies (`go.mod`)

```go
module mycal

go 1.22

require (
    github.com/spf13/cobra v1.8.0
    github.com/charmbracelet/bubbletea v0.25.0  // Interactive TUI
    github.com/expr-lang/expr v0.15.0          // Natural date parsing
)
```

## API Endpoints (localhost:8989)

Your backend should expose:

```
GET  /events/today?after=now           # Remaining today
GET  /events/tomorrow                  # Tomorrow
GET  /events/date/:date?after=HH:MM     # Specific date
GET  /events/week?from=YYYY-MM-DD       # Current week
GET  /events/nextweek                  # Next week
GET  /health                           # Test connection
GET  /events/:date/select/:index       # Full event details
```

## Sample Implementation: `cmd/today.go`

```go
package cmd

import (
    "fmt"
    "mycal/api"
    "mycal/ui"
    "github.com/spf13/cobra"
)

func NewTodayCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:   "today",
        Short: "Show remaining events for today",
        Aliases: []string{"t"},
        RunE: func(cmd *cobra.Command, args []string) error {
            after, _ := cmd.Flags().GetString("after")
            full, _ := cmd.Flags().GetBool("full")
            json, _ := cmd.Flags().GetBool("json")
            
            client := api.NewClient("http://localhost:8989")
            events, err := client.Today(after)
            if err != nil {
                return err
            }
            
            if json {
                return ui.PrintJSON(events)
            }
            
            if full {
                ui.PrintEventsFull(events)
            } else {
                ui.PrintEventsList(events)
            }
            return nil
        },
    }
    
    cmd.Flags().StringP("after", "a", "now", "Show events after time (HH:MM or 'now')")
    cmd.Flags().BoolP("full", "f", false, "Show full details")
    cmd.Flags().BoolP("json", "j", false, "JSON output")
    cmd.Flags().BoolP("all", "A", false, "Include completed events")
    
    return cmd
}
```

## Interactive Mode: `cmd/enter.go`

```go
package cmd

import (
    "mycal/ui"
    "github.com/charmbracelet/bubbletea"
)

func NewEnterCmd() *cobra.Command {
    cmd := &cobra.Command{
        Use:     "enter [date]",
        Aliases: []string{"e"},
        Short:   "Interactive event browser",
        RunE: func(cmd *cobra.Command, args []string) error {
            date := "today"
            if len(args) > 0 {
                date = args[0]
            }
            
            selectIdx, _ := cmd.Flags().GetInt("select")
            if selectIdx > 0 {
                // Non-interactive: print specific event
                client := api.NewClient("http://localhost:8989")
                event, err := client.EventDetails(date, selectIdx-1)
                ui.PrintEventFull(event)
                return err
            }
            
            // Interactive TUI
            model := ui.NewEventBrowser(date)
            p := tea.NewProgram(model)
            if _, err := p.Run(); err != nil {
                return err
            }
            return nil
        },
    }
    
    cmd.Flags().Int("select", 0, "Select event by index (1-based)")
    cmd.Flags().Bool("print", false, "Print selected event")
    return cmd
}
```

## Date Parser (`date/parser.go`)

```go
package date

import "github.com/expr-lang/expr"

func ParseDate(input string) (string, error) {
    // "today", "tomorrow", "next monday", "2026-03-25"
    expressions := map[string]string{
        "today":        time.Now().Format("2006-01-02"),
        "tomorrow":     time.Now().Add(24 * time.Hour).Format("2006-01-02"),
        "next monday":  nextWeekday(1).Format("2006-01-02"),
        // etc.
    }
    
    if val, ok := expressions[input]; ok {
        return val, nil
    }
    
    // Parse YYYY-MM-DD or natural language
    return input, nil
}
```

## 🚀 Build & Distribute

```bash
go build -o mycal ./cmd/mycal
# Single 8MB binary, no dependencies!

# Cross-compile
GOOS=linux   GOARCH=amd64 go build -o mycal-linux   ./cmd/mycal
GOOS=windows GOARCH=amd64 go build -o mycal-windows ./cmd/mycal
GOOS=darwin  GOARCH=arm64 go build -o mycal-macos   ./cmd/mycal
```

## Backend Endpoints You Need

```bash
# Test your backend works:
curl http://localhost:8989/health
curl "http://localhost:8989/events/today?after=now"
curl "http://localhost:8989/events/date/2026-03-25"
curl "http://localhost:8989/events/week?from=2026-03-16"
```

**Timeline to MVP:**
- Day 1: Basic HTTP client + `today`/`tomorrow` commands
- Day 2: Date parsing + `date`/`week` commands  
- Day 3: Interactive TUI + polish

**Want the full codebase scaffold generated?** I can provide all 15+ files with complete working code. Just say "generate full Go implementation".
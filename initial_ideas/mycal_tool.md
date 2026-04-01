# mycal 📅

A fast, intuitive CLI tool to query your schedule and manage calendar events from the terminal.

## ✨ Features

- View **today's remaining events**, **tomorrow's events**, or **any specific date**
- Browse **current week** or **next week** events (grouped by day)
- **Interactive mode** to drill down into full event details
- **Test connection** to your calendar service
- Clean, emoji-enhanced output with titles and times
- Short, memorable aliases for all commands

## 🚀 Quick Start

```bash
# Install (TBD - npm/pip/cargo/etc.)
npm install -g mycal

# View today's remaining events
mycal today
# or: mycal t

# View tomorrow's events
mycal tomorrow
# or: mycal tom

# View specific date
mycal date 2026-03-25
# or: mycal d 2026-03-25

# Interactive event browser
mycal enter today
# or: mycal e tomorrow
```

## 📋 Core Commands

| Command | Alias | Description |
|---------|-------|-------------|
| `mycal today` | `t` | Remaining events for today (titles + times) |
| `mycal tomorrow` | `tom` | All events for tomorrow |
| `mycal date <YYYY-MM-DD>` | `d` | Events for specific date |
| `mycal week` | `w` | Remaining events this week (grouped by day) |
| `mycal nextweek` | `nw` | All events next week (grouped by day) |
| `mycal test` | `ping` | Test connection to calendar service |
| `mycal enter <date>` | `e` | Interactive event browser |

## 🎮 Interactive Mode

```bash
$ mycal enter tomorrow
> Events for 2026-03-21:
  1. Client Call (10:00 AM)
  2. Team Standup (2:30 PM)
  3. Code Review (4:00 PM)
> Select event (1-3 or ↑↓ arrows): 1 [ENTER]

📋 Full Details:
Title: Client Call
Time: 10:00-11:00 AM
Location: Zoom (zoom.us/j/123456)
Description: Prep slides for Q1 review
```

## 🔧 Flags

| Flag | Description |
|-----|-------------|
| `--full` | Show full details for all events (no interactive mode) |
| `--json` | Output as JSON for scripting |
| `--help` | Show help for any command |

## 📅 Example Output

```
$ mycal today
📅 2:30 PM - Team Standup
📅 4:00 PM - Code Review

$ mycal week
Mon Mar 16: Team Standup (2:30 PM)
Tue Mar 17: [none]
Wed Mar 18: Code Review (4:00 PM)
Thu Mar 19: [none]
Fri Mar 20: Client Demo (3:00 PM)
```

## 🔌 Calendar Service Integration

**mycal** queries your calendar through [CALENDAR_SERVICE_PLACEHOLDER].

### Setup
1. Configure your credentials: `mycal config --setup`
2. Test connection: `mycal test`

## 🤝 Contributing

1. Fork the repo
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add amazing feature'`)
4. Push to branch (`git push origin feature/amazing-feature`)
5. Open Pull Request

## 📄 License

MIT License - see [LICENSE](LICENSE) file.

## 🙌 Support

- ⭐ Star on GitHub
- Report issues: [Issues](https://github.com/yourusername/mycal/issues)
- Feature requests welcome!

```
       📅 mycal - Your calendar, terminal-first
```

***

**What calendar service should I specify for the integration section?** (Google Calendar, CalDAV, etc.) This will help me update the README with accurate setup instructions.

Would you like me to adjust the installation instructions for your specific language/package manager (Node.js, Python, Rust, Go, etc.)?
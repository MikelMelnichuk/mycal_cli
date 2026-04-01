## Here is a versatile command usage of mycal tool

### Today (Remaining Events)
```
# Basic: Shows 2 upcoming events
$ mycal today
📅 3:00 PM - Client Demo
📅 5:30 PM - 1:1 with Boss

# With full details flag
$ mycal today --full
📅 3:00-3:45 PM - Client Demo (Zoom: zoom.us/j/123, Notes: Show v2 features)

# JSON for tests/parsing
$ mycal today --json
[
  {"title": "Client Demo", "start": "2026-03-20T15:00:00", "end": "2026-03-20T15:45:00"},
  {"title": "1:1 with Boss", "start": "2026-03-20T17:30:00"}
]

# Edge: No events left (after 6 PM)
$ mycal today  # => "No events remaining today."
```

### Tomorrow
```
# Basic
$ mycal tomorrow
📅 10:00 AM - Team Sync
📅 2:00 PM - Lunch w/ Team

# Alias + filter by time
$ mycal tom --after 12:00
📅 2:00 PM - Lunch w/ Team

# Weekend edge (no events)
$ mycal tomorrow  # => "No events tomorrow (Saturday)."
```

### Date
```
# Specific date
$ mycal date 2026-03-25
📅 9:00 AM - Dentist
📅 7:00 PM - Movie Night

# Natural language (if supported)
$ mycal date next monday
$ mycal date "this friday"
$ mycal date "Mar 25"

# Past date (shows all, even completed)
$ mycal date 2026-03-18 --completed
```

### Week (Current Week, Remaining)
```
# Basic (groups by day, remaining only)
$ mycal week
Mon Mar 16: [completed - skipped]
Tue Mar 17: [completed]
Wed Mar 18: Code Review (4:00 PM) ✅
Thu Mar 19: [none]
Fri Mar 20: Client Demo (3:00 PM), 1:1 (5:30 PM)
Sat/Sun: [weekend]

# Full week including past
$ mycal week --all
$ mycal w --json | jq '.[] | select(.day == "Fri")'
```

### Next Week
```
# Basic
$ mycal nextweek
Mar 23-29:
Mon: Project Kickoff (11:00 AM)
Tue: [none]
...

```

### Test Connection
```
# Basic ping
$ mycal test
✅ Connected to calendar service. Last sync: 2026-03-20 14:32 MSK
Latency: 120ms

# Verbose with auth check
$ mycal test --verbose
✅ Auth: Valid (expires 2026-04-20)
✅ API: Google Calendar v3 reachable
❌ Sync lag: 5min behind

# Fail cases for tests
$ mycal test --json  # => {"status": "error", "message": "Auth expired"}
```

### Enter (Interactive Event Details)
```
# Start interactive for today
$ mycal enter today
Loaded 2 events for Fri Mar 20:
1. Client Demo (3:00 PM)
2. 1:1 with Boss (5:30 PM)
> 1 [ENTER]
---
Event: Client Demo
When: 2026-03-20 15:00-15:45 MSK
Where: Zoom (zoom.us/j/123)
Attendees: client@ex.com, you@ex.com
Description: Demo v2 features. Prep slides.
[Press Q to quit or ↑↓ to select]

# Direct non-interactive
$ mycal enter today --select 1 --print
# Outputs full JSON/event details

# Edge: Empty day
$ mycal enter tomorrow  # => "No events. Press any key to exit."
```

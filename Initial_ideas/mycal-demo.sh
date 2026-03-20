#!/bin/bash
# mycal usage examples

echo '$ mycal today'
echo '📅 3:00 PM - Client Demo'
echo '📅 5:30 PM - 1:1 with Boss'
echo

echo '$ mycal today --full'
echo '📅 3:00-3:45 PM - Client Demo (Zoom: zoom.us/j/123, Notes: Show v2 features)'
echo

echo '$ mycal today --json'
echo '[{"title": "Client Demo", "start": "2026-03-20T15:00:00", "end": "2026-03-20T15:45:00"},'
echo ' {"title": "1:1 with Boss", "start": "2026-03-20T17:30:00"}]'
echo

echo '$ mycal today'
echo 'No events remaining today.'
echo

echo '$ mycal tomorrow'
echo '📅 10:00 AM - Team Sync'
echo '📅 2:00 PM - Lunch w/ Team'
echo

echo '$ mycal tom --after 12:00'
echo '📅 2:00 PM - Lunch w/ Team'
echo

echo '$ mycal tomorrow'
echo 'No events tomorrow (Saturday).'
echo

echo '$ mycal date 2026-03-25'
echo '📅 9:00 AM - Dentist'
echo '📅 7:00 PM - Movie Night'
echo

echo '$ mycal date next monday'
echo '📅 11:00 AM - Project Kickoff'
echo

echo '$ mycal date 2026-03-18 --completed'
echo '📅 2:00 PM - Past Meeting ✅'
echo

echo '$ mycal week'
echo 'Mon Mar 16: [completed - skipped]'
echo 'Tue Mar 17: [completed]'
echo 'Wed Mar 18: Code Review (4:00 PM) ✅'
echo 'Thu Mar 19: [none]'
echo 'Fri Mar 20: Client Demo (3:00 PM), 1:1 (5:30 PM)'
echo 'Sat/Sun: [weekend]'
echo

echo '$ mycal w --all'
echo 'Full week view including completed events...'
echo

echo '$ mycal w --json | jq ".[] | select(.day == \"Fri\")"'
echo '[{"day": "Fri", "events": ["Client Demo 3PM", "1:1 5:30PM"]}]'
echo

echo '$ mycal nextweek'
echo 'Mar 23-29:'
echo 'Mon: Project Kickoff (11:00 AM)'
echo 'Tue: [none]'
echo 'Wed: Offsite Meeting (9:00 AM)'
echo 'Thu: [none]'
echo 'Fri: Team Happy Hour (6:00 PM)'
echo

echo '$ mycal nw --compact'
echo 'Mon: Project Kickoff 11AM | Wed: Offsite 9AM | Fri: Happy Hour 6PM'
echo

echo '$ mycal test'
echo '✅ Connected to calendar service. Last sync: 2026-03-20 14:32 MSK'
echo 'Latency: 120ms'
echo

echo '$ mycal test --verbose'
echo '✅ Auth: Valid (expires 2026-04-20)'
echo '✅ API: Google Calendar v3 reachable'
echo '✅ Sync: Up to date'
echo

echo '$ mycal test --json'
echo '{"status": "success", "latency": 120, "last_sync": "2026-03-20T14:32:00"}'
echo

echo '$ mycal enter today'
echo 'Loaded 2 events for Fri Mar 20:'
echo '1. Client Demo (3:00 PM)'
echo '2. 1:1 with Boss (5:30 PM)'
echo '> 1 [ENTER]'
echo '---'
echo 'Event: Client Demo'
echo 'When: 2026-03-20 15:00-15:45 MSK'
echo 'Where: Zoom (zoom.us/j/123)'
echo 'Attendees: client@ex.com, you@ex.com'
echo 'Description: Demo v2 features. Prep slides.'
echo '[Press Q to quit or ↑↓ to select]'
echo

echo '$ mycal enter today --select 1 --print'
echo '{"title": "Client Demo", "details": "Full event JSON output..."}'
echo

echo '$ mycal enter tomorrow'
echo 'No events. Press any key to exit.'
echo

echo '$ mycal t'
echo '📅 3:00 PM - Client Demo'
echo

echo '$ mycal tom'
echo '📅 10:00 AM - Team Sync'
echo

echo '$ mycal d 2026-03-25'
echo '📅 9:00 AM - Dentist'
echo

echo '$ mycal w'
echo 'Fri Mar 20: 2 events remaining'
echo

echo '$ mycal nw'
echo 'Next week: 3 events'
echo

echo '$ mycal ping'
echo '✅ Connection OK'

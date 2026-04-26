mycal today
mycal today --all
mycal today --json
mycal today --after 13:00
mycal t

mycal tomorrow
mycal tomorrow --json
mycal tomorrow --after 12:00
mycal tom

mycal date YYYY-MM-DD
mycal date 2027-03-25 # Date in the future
mycal date 2026-03-18 # Date in the past
mycal date YYYY-MM-DD --json
mycal date next monday
mycal d

mycal week       # Shows the events that are yet to happen (or in progress), without the past events
mycal week --all # Shows all the events events of the current weak (including past)
mycal week --json | jq ".[] | select(.day == \"Fri\")"
mycal w

mycal nextweek
mycal nextweek --json
mycal nw

# This option tests the connection to the server
mycal test
mycal test --verbose

# This option helps the user inspect what events he has on a given date, and check in more detail their content
mycal enter today
mycal enter today --select 1 --print
mycal enter tomorrow
mycal enter next friday
mycal enter YYYY-MM-DD

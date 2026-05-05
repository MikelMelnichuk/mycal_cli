package models

type Event struct {
	ID    string `json:"id"`
	Title string `json:"title"`
	Start string `json:"start"`
	End   string `json:"end"`
	Day   string `json:"day"` // YYYY-MM-DD format (e.g., "2026-03-18")
}

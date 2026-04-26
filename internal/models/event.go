package models

import "time"

type Event struct {
	ID    string    `json:"id"`
	Title string    `json:"title"`
	Start time.Time `json:"start"`
	End   time.Time `json:"end"`
	Day   string    `json:"day"` // e.g., "2026-03-18"
}

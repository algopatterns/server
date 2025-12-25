package examples

import (
	"time"
)

// internal enriched representation
type Example struct {
	ID          string
	Title       string
	Description string
	Code        string
	Tags        []string
	Author      string
	Category    string
	SourceURL   string
	CreatedAt   time.Time
}

// represents input from external sources (JSON file, API request, etc.)
type RawExample struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Code        string   `json:"code"`
	Tags        []string `json:"tags"`
	Category    string   `json:"category"`
	Author      string   `json:"author"`
}

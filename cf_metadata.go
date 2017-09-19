package main

import "time"

// CFMetadata is a set of data any Cloud Foundry entity has.
type CFMetadata struct {
	GUID      string    `json:"guid"`
	URL       string    `json:"url"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

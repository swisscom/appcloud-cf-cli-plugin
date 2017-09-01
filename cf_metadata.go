package main

// CFMetadata is a set of data any Cloud Foundry entity has.
type CFMetadata struct {
	GUID      string `json:"guid"`
	URL       string `json:"url"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
}

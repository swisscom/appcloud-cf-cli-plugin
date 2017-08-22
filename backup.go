package main

// BackupsResponse is the response of the server to a get backups call
type BackupsResponse struct {
	TotalResult int      `json:"total_results"`
	TotalPages  int      `json:"total_pages"`
	PrevURL     string   `json:"prev_url"`
	NextURL     string   `json:"next_url"`
	Resources   []Backup `json:"resources"`
}

// Backup is a service instance backup
type Backup struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		ServiceInstanceID string `json:"service_instance_id"`
		Status            string `json:"status"`
	} `json:"entity"`
}

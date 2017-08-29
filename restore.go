package main

// Restore is a service instance restore
type Restore struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		BackupID string `json:"service_instance_id"`
		Status            string `json:"status"`
	} `json:"entity"`
}

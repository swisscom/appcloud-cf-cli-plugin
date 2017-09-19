package main

// RestoreBackupResponse is a response from the server to a restore backup call.
type RestoreBackupResponse struct {
	Restore
	ServerResponseError
}

// Restore is a service instance restore
type Restore struct {
	Metadata CFMetadata `json:"metadata"`
	Entity   struct {
		BackupID string `json:"service_instance_id"`
		Status   string `json:"status"`
	} `json:"entity"`
}

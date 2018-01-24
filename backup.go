package main

// Backup is a service instance backup.
type Backup struct {
	Metadata CFMetadata `json:"metadata"`
	Entity   struct {
		Status   string    `json:"status"`
		Restores []Restore `json:"restores"`
	} `json:"entity"`
}

// BackupResponse is the response from the server from a create backup call.
type BackupResponse struct {
	Backup
	ServerResponseError
}

// BackupsResponse is the response of the server to a get backups call.
type BackupsResponse struct {
	Resources []Backup `json:"resources"`
	ServerResponsePagination
	ServerResponseError
}

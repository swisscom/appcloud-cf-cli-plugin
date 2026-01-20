package appcloud

// EventsResponse is a response from the server to an events call.
type EventsResponse struct {
	Resources []V3AuditEvent `json:"resources"`
	// ServerResponsePagination
	V3ServerResponseErrors
}

// V3AuditEvent is a service instance event.
type V3AuditEvent struct {
	CreatedAt string `json:"created_at"`
	Type      string `json:"type"`
	Actor     struct {
		Name string `json:"name"`
	} `json:"actor"`
}

package main

// EventsResponse is a response from the server to an events call.
type EventsResponse struct {
	Resources []Event `json:"resources"`
	ServerResponsePagination
	ServerResponseError
}

// Event is a service instance event.
type Event struct {
	Metadata CFMetadata `json:"metadata"`
	Entity   struct {
		Type      string `json:"type"`
		ActorName string `json:"actor_name"`
	} `json:"entity"`
}

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
		Actor     string `json:"actor"`
		ActorType string `json:"actor_type"`
		ActorName string `json:"actor_name"`
		Actee     string `json:"actee"`
		ActeeType string `json:"actee_type"`
		ActeeName string `json:"actee_name"`
		TimeStamp string `json:"timestamp"`
		SpaceGUID string `json:"space_guid"`
		OrgGUID   string `json:"organization_guid"`
	} `json:"entity"`
}

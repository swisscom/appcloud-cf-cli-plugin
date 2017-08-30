package main

// ServiceEvent provides service instance events
type ServiceEvent struct {
	TotalResults int               `json:"total_results"`
	TotalPages   int               `json:"total_pages"`
	PrevUrl      string            `json:"prev_url"`
	NextUrl      string            `json:"next_url"`
	Resources    []ServiceInstanceEventData `json:"resources"`
}

type ServiceInstanceEventData struct {
		Metadata struct {
			GUID      string `json:"guid"`
			URL       string `json:"url"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"metadata"`
		Entity struct {
			Type string `json:"type"`
			Actor string `json:"actor"`
			ActorType string `json:"actor_type"`
			ActorName string `json:"actor_name"`
			ActorUserName string `json:"actor_username"`
			Actee string `json:"actee"`
			ActeeType string `json:"actee_type"`
			ActeeName string `json:"actee_name"`
			TimeStamp string `json:"timestamp"`
			SpaceGUID string `json:"space_guid"`
			OrgGUID string `json:"organization_guid"`
			ServiceMetadata struct {
				Request struct {
					Name 			string `json:"maridadb"`
					SpaceGUID  		string `json:"space_guid"`
					ServicePlanGUID string `json:"service_plan_guid"`
					Parameters		string `json:"parameters"`
					Tags			string `json:"tags"`
				}`json:"metadata"`
			} `json:"metadata"`
		} `json:"entity"`
	}

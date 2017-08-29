package main

// Backup is a service instance backup
type OrgInvitation struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Invitee          string `json:"invitee"`
		OrgID            string `json:"organization_id"`
		ActorUserName    string `json:"actor_username"`
		ActorUsreID      string `json:"actor_user_id"`
		OrganizationName string `json:"organization_name"`
		AvatarUrl        string `json:"avatar_url"`
		Status           string `json:"status"`
	} `json:"entity"`
}

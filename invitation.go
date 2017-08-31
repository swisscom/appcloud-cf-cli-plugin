package main

// Invitation is an invitation a user received to join a specific entity.
type Invitation struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Invitee          string   `json:"invitee"`
		Roles            []string `json:"roles"`
		ActorUsername    string   `json:"actor_username"`
		ActorUserID      string   `json:"actor_user_id"`
		AccountID        string   `json:"account_id"`
		AccountName      string   `json:"account_name"`
		OrganizationID   string   `json:"organization_id"`
		OrganizationName string   `json:"organization_name"`
		SpaceID          string   `json:"space_id"`
		SpaceName        string   `json:"space_name"`
		Status           string   `json:"status"`
	} `json:"entity"`
}

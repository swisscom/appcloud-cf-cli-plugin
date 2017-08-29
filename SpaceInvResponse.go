package main

type SpaceInvResponse struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Invitee          string   `json:"invitee"`
		Roles            []string `json:"roles"`
		SpaceID          string   `json:"space_id"`
		ActorUserName    string   `json:"actor_username"`
		ActorUserId      string   `json:"actor_user_id"`
		SpaceName        string   `json:"space_name"`
		OrganizationId   string   `json:"organization_id"`
		OrganizationName string   `json:"organization_name"`
		AvatarUrl        string   `json:"avatar_url"`
		Status           string   `json:"status"`
	} `json:"entity"`
}

package main

type OrganizationInvResponse struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		Invitee string `json:"invitee"`
        Roles [] string `json:"roles"`
		ActorUserName string `json:"actor_username"`
		ActorUserId string `json:"actor_user_id"`
        OrganizationId string `json:"organization_id"`
        OrganizationName string `json:"organization_name"`
        AvatarUrl string `json:"avatar_url"`
        Status string `json:"status"`		
	} `json:"entity"`
}
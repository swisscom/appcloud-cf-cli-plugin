package main

//Tree Structs
type OrgResponse struct {
	GenericResponse
	Resources []Organization `json:"resources"`
}

type Application struct {
	ID string `json:"id"`
	Name string `json:"name"`
	BackupIconUrl string `json:"buildpack_icon_url"`
}

type ServiceInstance struct {
	ID string `json:"id"`
	Name string `json:"name"`
	ServiceIconUrl string `json:"service_icon_url"`
}

type Space struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Applications []Application `json:"applications"`
	ServiceInstances []ServiceInstance `json:"service_instances"`
}

type Organization struct {
	ID string `json:"id"`
	Name string `json:"name"`
	Status string `json:"status"`
	Spaces []Space `json:"spaces"`
}

//Invitation Structs
type OrganizationEntity struct {
	InvitationEntity
	OrganizationId string `json:"organization_id"`
	OrganizationName string `json:"invitee"`
}

type InvitationEntity struct {
	Invitee string `json:"invitee"`
	Roles []string `json:"roles"`
	ActorUsername string `json:"actor_username"`
	ActorUserId string `json:"actor_user_id"`
	AvatarUrl string `json:"avatar_url"`
	Status string `json:"status"`
}

type SpaceInvitation struct {
	Metadata struct {
			 GUID string `json:"guid"`
			 URL string `json:"url"`
			 CreatedAt string `json:"created_at"`
			 UpdatedAt string `json:"updated_at"`
		 } `json:"metadata"`
	SpaceEntity struct {
			    OrganizationEntity
			    SpaceId string `json:"space_id"`
			    SpaceName string `json:"space_name"`
		    } `json:"entity"`
}

type OrganizationInvitation struct {
	Metadata struct {
			 GUID string `json:"guid"`
			 URL string `json:"url"`
			 CreatedAt string `json:"created_at"`
			 UpdatedAt string `json:"updated_at"`
		 } `json:"metadata"`
	OrganizationEntity struct {
				   InvitationEntity
				   OrganizationId string `json:"organization_id"`
				   OrganizationName string `json:"organization_name"`
			   } `json:"entity"`
}

type AccountInvitation struct {
	Metadata struct {
		 GUID string `json:"guid"`
		 URL string `json:"url"`
		 CreatedAt string `json:"created_at"`
		 UpdatedAt string `json:"updated_at"`
	 } `json:"metadata"`
	AccountEntity struct {
	         InvitationEntity
	         AccountId string `json:"account_id"`
	         AccountName string `json:"account_name"`
         } `json:"entity"`
}

type SpaceResponse struct {
	GenericResponse
	Resources []SpaceInvitation `json:"resources"`
}

type OrganizationResponse struct {
	GenericResponse
	Resources []OrganizationInvitation `json:"resources"`
}

type AccountResponse struct {
	GenericResponse
	Resources []AccountInvitation `json:"resources"`
}

type GenericResponse struct {
	TotalResults int `json:"total_results"`
	TotalPages int `json:"total_pages"`
	PrevUrl string `json:"prev_url"`
	NextUrl string `json:"next_url"`
}

type DockerRepository struct {
	Repositories []string `json:"repositories"`
}
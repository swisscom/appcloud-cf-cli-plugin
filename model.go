package main

type OrgResponse struct {
	TotalResults int `json:"total_results"`
	TotalPages int `json:"total_pages"`
	PrevUrl string `json:"prev_url"`
	NextUrl int `json:"next_url"`
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

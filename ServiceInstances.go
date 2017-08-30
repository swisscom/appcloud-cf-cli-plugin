package main

// ServiceInstance provides service instance details
type SpaceSummary struct {
	Name 	 string      `json:"name"`
	GUID   	 string      `json:"guid"`
	Apps     []App       `json:"apps"`
	Services []Service 	`json:"services"`
}
type App struct{
		Name      			string `json:"name"`
		Memory 				int `json:"memory"`
		Instances 			int `json:"instances"`
		State 				string `json:"state"`
		DetectedBuildpack 	string `json:"detected_buildpack"`
		Buildpack 			string `json:"buildpack"`
		DockerImage 		string `json:"docker_image"`
		DockerCredentials   string `json:"docker_credentials"`
		BuildpackIconURL 	string `json:"buildpack_icon_url"`
		Urls 				[]string `json:"urls"`
		UrlsWithSceme 		[]string `json:"urls_with_scheme"`
		SpaceGUID 			string `json:"space_guid"`
		GUID 				string `json:"guid"`
		Routes 				[]Route `json:"routes"`
		ServiceCount      	int `json:"service_count"`
		ServiceNames 		[]string `json:"service_names"`
		RunningInstances 	int `json:"running_instances"`
		ServicePlan 		ServicePlan `json:"service_plan"`	
		BoundAppsCount 		int `json:"bound_app_count"`
		Production 			bool `json:"production"`
		StackGUID      		string `json:"stack_guid"`
		EnvJson struct{
			key string
		}`json:"environment_json"`
		DiskQuota 			int `json:"disk_quota"`
		Version 			string `json:"version"`
		Command 			string `json:"command"`
		Console 			bool `json:"console"`
		Debug      			string `json:"debug"`
		StagingTaskID 		string `json:"staging_task_id"`
		PackageState 		string `json:"package_state"`
		HealthCheckType 	string `json:"health_check_type"`
		HealthCheckTimeout 			int `json:"health_check_timeout"`
		HealthCheckHttpCheckpoint 	int `json:"health_check_http_endpoint"`
		StagingFailedReason string `json:"staging_failed_reason"`
		StagingFailedDesc   string `json:"staging_failed_description"`
		Diego 				bool `json:"diego"`
		PkgUpdatedAt 		string `json:"package_updated_at"`
		DetectedStartCmd 	string `json:"detected_start_command"`
		EnableSSH 			bool `json:"enable_ssh"`
		Ports 				string `json:"ports"`
}


type Route struct{
		GUID 	string `json:"guid"`
		Host 	string `json:"host"`
		Port 	string `json:"port"`
		Path 	string `json:"path"`
		Domain struct{
			GUID 	string `json:"guid"`
			Name 	string `json:"name"`
		}`json:"domain"`
}

type Service struct{
		Name      		string `json:"name"`
		ServiceKeyCount string `json:"service_key_count"`
		BoundRouteCount int `json:"bound_route_count"`
		ServicePlan 	ServicePlan `json:"service_plan"`
		GUID 			string `json:"guid"`
		BoundAppsCount 	int `json:"bound_app_count"`
		DashBoardURL 	string `json:"dashboard_url"`
		LastOperation 	struct{
			Type			string `json:"type"`
			State      		string `json:"state"`
			Description     string `json:"description"`
			UpdatedAt      	string `json:"updated_at"`
			CreatedAt      	string `json:"created_at"`
		} `json:"last_operation"`
}

type ServicePlan struct{
		Name      	string `json:"name"`
		Service 	struct {
			Label 				string `json:"label"`
			ServiceIconURL      string `json:"service_icon_url"`	
			Version            	string `json:"version"`
			Provider            string `json:"provider"`
			GUID 	string `json:"guid"`
		} `json:"service"`
	Guid string `json:"guid"`
}

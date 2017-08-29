package main

// SSL Certificate
type SSLCertificate struct {
	Metadata struct {
		GUID      string `json:"guid"`
		URL       string `json:"url"`
		CreatedAt string `json:"created_at"`
		UpdatedAt string `json:"updated_at"`
	} `json:"metadata"`
	Entity struct {
		SpaceID        string `json:"space_id"`
		CertificateID  string `json:"certificate_id"`
		FullDomainName string `json:"full_domain_name"`
		ProcessType    string `json:"process_type"`
		CurrentStatus  string `json:"current_status"`
	} `json:"entity"`
}

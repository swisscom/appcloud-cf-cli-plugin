package main

// SSL Certificate
type SSLCertificates struct {
	TotalResults int `json:"total_results"`
	TotalPages int `json:"total_pages"`
	PrevUrl string `json:"prev_url"`
	NextUrl string `json:"next_url"`
	Resources []CertificateData `json:"resources"`
}

type CertificateData struct {
		Metadata struct {
			GUID      string `json:"guid"`
			URL       string `json:"url"`
			CreatedAt string `json:"created_at"`
			UpdatedAt string `json:"updated_at"`
		} `json:"metadata"`
		Entity struct {
			SpaceID string `json:"space_id"`
			Status string `json:"status"`
			FullDomainName string `json:"full_domain_name"`
			NotValidBefore string `json:"not_valid_before"`
			NotValidAfter string `json:"not_valid_after"`
			AutomaticRenewal string `json:"automatic_renewal"`
		} `json:"entity"`
}

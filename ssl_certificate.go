package main

// SSLCertificate is an SSL certificate which can be used to secure routes.
type SSLCertificate struct {
	Metadata CFMetadata `json:"metadata"`
	Entity   struct {
		FullDomainName string `json:"full_domain_name"`
		Status         string `json:"status"`
	} `json:"entity"`
}

// SSLCertificateRequest is the request to create an SSL certificate.
type SSLCertificateRequest struct {
	SpaceID        string `json:"space_id"`
	FullDomainName string `json:"full_domain_name"`
}

// SSLCertificateResponse is the response from the server from a create certificate call
type SSLCertificateResponse struct {
	SSLCertificate
	ServerResponseError
}

// SSLCertificatesResponse is a response from the server .
type SSLCertificatesResponse struct {
	Resources []SSLCertificate `json:"resources"`
	ServerResponsePagination
	ServerResponseError
}

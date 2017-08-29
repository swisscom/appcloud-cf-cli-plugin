package main

// ServiceBrokerResponse is the response from a service broker potentially containing errors
type ServerResponse struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	ErrorCode   string `json:"error_code"`
}

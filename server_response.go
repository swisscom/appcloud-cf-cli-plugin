package main

// ServerResponseError is the response from a server which may contain an error.
type ServerResponseError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	ErrorCode   string `json:"error_code"`
}

// ServerResponsePagination is a server response including pagination information.
type ServerResponsePagination struct {
	TotalResult int    `json:"total_results"`
	TotalPages  int    `json:"total_pages"`
	PrevURL     string `json:"prev_url"`
	NextURL     string `json:"next_url"`
}

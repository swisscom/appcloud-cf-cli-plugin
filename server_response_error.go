package main

// ServerResponseError is the response from a server which may contain an error
type ServerResponseError struct {
	Code        int    `json:"code"`
	Description string `json:"description"`
	ErrorCode   string `json:"error_code"`
}

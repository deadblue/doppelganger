package json

import "encoding/json"

type Request struct {
	// Request ID
	Id string `json:"id"`
	// Method name
	Method string `json:"method"`
	// Parameters
	Params json.RawMessage `json:"params"`
}

type Response struct {
	// Request ID
	Id string `json:"id"`
	// Error code
	Error int `json:"error"`
	// Error message when error code is non-zero
	Message string `json:"message,omitempty"`
	// Result when error code is zero.
	Result interface{} `json:"result,omitempty"`
}

package json

type Response struct {
	// Error code
	Error int `json:"error"`
	// Error message when error code is non-zero
	Message string `json:"message,omitempty"`
	// Result when error code is zero.
	Result interface{} `json:"result,omitempty"`
}

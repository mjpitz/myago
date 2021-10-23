package yarpc

// Status reports an optional code and message along with the request.
type Status struct {
	Code    int    `json:"code,omitempty"`
	Message string `json:"message,omitempty"`
}

// Frame is the generalized structure passed along the wire.
type Frame struct {
	Nonce  string      `json:"nonce,omitempty"`
	Status *Status     `json:"status,omitempty"`
	Body   interface{} `json:"body"`
}

type Invoke struct {
	Method string `json:"method,omitempty"`
}

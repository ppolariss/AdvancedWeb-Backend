package message

type AIResponse struct {
	Status     int    `json:"status"` // 1 for output, 0 for end, -1 for error
	StatusCode int    `json:"status_code,omitempty"`
	Output     string `json:"output,omitempty"`
	Stage      string `json:"stage,omitempty"`
}

type MossResponse struct {
	Status int    `json:"status"`
	Output string `json:"output"`
	Stage  string `json:"stage"`
}

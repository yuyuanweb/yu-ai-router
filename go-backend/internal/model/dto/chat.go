package dto

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model       string        `json:"model"`
	Messages    []ChatMessage `json:"messages"`
	Stream      *bool         `json:"stream"`
	Temperature *float64      `json:"temperature"`
	MaxTokens   *int          `json:"max_tokens"`
}

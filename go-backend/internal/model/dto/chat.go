package dto

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatRequest struct {
	Model           string        `json:"model"`
	Messages        []ChatMessage `json:"messages"`
	Stream          *bool         `json:"stream"`
	Temperature     *float64      `json:"temperature"`
	MaxTokens       *int          `json:"max_tokens"`
	EnableReasoning *bool         `json:"enable_reasoning"`
	RoutingStrategy string        `json:"routing_strategy"`
	PluginKey       string        `json:"plugin_key"`
	FileURL         string        `json:"file_url"`
	FileBytes       []byte        `json:"-"`
	FileType        string        `json:"file_type"`
}

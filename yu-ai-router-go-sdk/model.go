package yuairoutersdk

type ChatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

func UserMessage(content string) ChatMessage {
	return ChatMessage{Role: "user", Content: content}
}

func SystemMessage(content string) ChatMessage {
	return ChatMessage{Role: "system", Content: content}
}

func AssistantMessage(content string) ChatMessage {
	return ChatMessage{Role: "assistant", Content: content}
}

type ChatRequest struct {
	Model           string        `json:"model,omitempty"`
	Messages        []ChatMessage `json:"messages"`
	Stream          *bool         `json:"stream,omitempty"`
	Temperature     *float64      `json:"temperature,omitempty"`
	MaxTokens       *int          `json:"max_tokens,omitempty"`
	EnableReasoning *bool         `json:"enable_reasoning,omitempty"`
	RoutingStrategy string        `json:"routing_strategy,omitempty"`
}

func SimpleRequest(message string) ChatRequest {
	return ChatRequest{
		Messages: []ChatMessage{
			UserMessage(message),
		},
	}
}

func RequestWithModel(model, message string) ChatRequest {
	request := SimpleRequest(message)
	request.Model = model
	return request
}

type ChatResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Model   string               `json:"model"`
	Choices []ChatResponseChoice `json:"choices"`
	Usage   ChatResponseUsage    `json:"usage"`
}

type ChatResponseChoice struct {
	Index        int         `json:"index"`
	Message      ChatMessage `json:"message"`
	FinishReason string      `json:"finishReason"`
}

type ChatResponseUsage struct {
	PromptTokens     int `json:"promptTokens"`
	CompletionTokens int `json:"completionTokens"`
	TotalTokens      int `json:"totalTokens"`
}

type StreamResponse struct {
	ID      string                 `json:"id"`
	Object  string                 `json:"object"`
	Created int64                  `json:"created"`
	Model   string                 `json:"model"`
	Choices []StreamResponseChoice `json:"choices"`
}

type StreamResponseChoice struct {
	Index        int                 `json:"index"`
	Delta        StreamResponseDelta `json:"delta"`
	FinishReason string              `json:"finishReason,omitempty"`
}

type StreamResponseDelta struct {
	Role             string `json:"role,omitempty"`
	Content          string `json:"content,omitempty"`
	ReasoningContent string `json:"reasoningContent,omitempty"`
}

type ChatChunk struct {
	Content          string
	ReasoningContent string
	Model            string
	Done             bool
}

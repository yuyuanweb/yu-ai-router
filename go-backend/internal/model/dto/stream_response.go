package dto

type StreamResponse struct {
	ID      string               `json:"id"`
	Object  string               `json:"object"`
	Created int64                `json:"created"`
	Model   string               `json:"model"`
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

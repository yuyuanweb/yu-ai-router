package dto

type ImageGenerationRequest struct {
	Prompt         string `json:"prompt"`
	Model          string `json:"model"`
	N              *int   `json:"n"`
	Size           string `json:"size"`
	Quality        string `json:"quality"`
	ResponseFormat string `json:"response_format"`
	User           string `json:"user"`
}

type ImageGenerationResponse struct {
	Created int64            `json:"created"`
	Data    []ImageDataEntry `json:"data"`
}

type ImageDataEntry struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revisedPrompt,omitempty"`
}
